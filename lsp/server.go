// Package lsp is a partial Language Server Protocol implementation. WhaleLint provides it's findings through
// Diagnostics.
package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	Log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	Linter "github.com/cremindes/whalelint/linter"
	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

type (
	MethodMapType          map[string]func(interface{}) (interface{}, error)
	NotificationHandlerMap map[string]func([]byte) (interface{}, error)
)

var (
	MethodMap              MethodMapType          // nolint:gochecknoglobals
	notificationHandlerMap NotificationHandlerMap // nolint:gochecknoglobals
)

type TextDocumentURIandStageList struct {
	StageList []instructions.Stage
	URI       DocumentURI
}

// Yay is a dummy function for notifications that are not yet supported or we do not care about them.
func Yay(_ []byte) (interface{}, error) {
	Log.Println("Yay")

	return nil, nil
}

func OnTextOpen(requestBytes []byte) (interface{}, error) {
	type TextDocumentWrapper struct {
		TextDocument TextDocumentItem `json:"textDocument"`
	}

	testDocParam := struct {
		JSONrpcVersion string              `json:"jsonrpc"`      // "jsonrpc": "2.0",
		ID             interface{}         `json:"id,omitempty"` // "id": 1,
		Method         string              `json:"method"`       // "method": "textDocument/didOpen",
		Params         TextDocumentWrapper `json:"params"`
	}{}

	err := json.Unmarshal(requestBytes, &testDocParam)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONRPC request OnTextOpen: %w", err)
	}

	stageList := parseFromText(testDocParam.Params.TextDocument.Text)

	result := TextDocumentURIandStageList{
		StageList: stageList,
		URI:       testDocParam.Params.TextDocument.URI,
	}

	return result, nil
}

func onTextDocumentDidChange(requestBytes []byte) (interface{}, error) {
	testDocParam := struct {
		JSONrpcVersion string                      `json:"jsonrpc"`      // "jsonrpc": "2.0",
		ID             interface{}                 `json:"id,omitempty"` // "id": 1,
		Method         string                      `json:"method"`       // "method": "textDocument/didChange",
		Params         DidChangeTextDocumentParams `json:"params"`
	}{}

	err := json.Unmarshal(requestBytes, &testDocParam)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSONRPC request OnTextOpen: %w", err)
	}

	stageList := parseFromText(testDocParam.Params.ContentChanges[0].Text)

	result := TextDocumentURIandStageList{
		StageList: stageList,
		URI:       testDocParam.Params.TextDocument.URI,
	}

	return result, nil
}

func onTextDocumentDidSave(requestBytes []byte) (interface{}, error) {
	Log.Println("Save")

	return "", nil
}

func parseFromText(str string) []instructions.Stage {
	// update RawParser
	Parser.RawParser.UpdateRawStr(str)

	reader := strings.NewReader(str)

	dockerfile, err := parser.Parse(reader)
	if err != nil {
		Log.Error("Cannot parse Dockerfile", err)
	}

	stageList, _ := Utils.ParseDockerfileInstructionsSafely(dockerfile, reader)

	return stageList
}

func PublishDiagnostics(uriAndStageList TextDocumentURIandStageList, w *bufio.Writer) {
	rr := PublishDiagnosticsParams{
		URI:         uriAndStageList.URI,
		Version:     0,
		Diagnostics: nil,
	}

	// lint
	diagList := Linter.MainLinter.Run(uriAndStageList.StageList)
	violationList := filter.Choose(diagList,
		func(x RuleSet.RuleValidationResult) bool {
			return x.IsViolated()
		}).([]RuleSet.RuleValidationResult)

	rr.Diagnostics = make([]Diagnostic, len(violationList))

	for i, diag := range violationList {
		rr.Diagnostics[i] = Diagnostic{
			Range:              VSCodeRangeFromLocationRange(diag.LocationRange),
			Severity:           VSCodeSeverityFromSeverity(diag.Severity()),
			Code:               diag.RuleID(),
			CodeDescription:    nil,
			Source:             "WhaleLint",
			Message:            diag.Message(),
			Tags:               nil,
			RelatedInformation: nil,
			Data:               nil,
		}
	}

	rpcResponse := &RPCNotification{
		Method: "textDocument/publishDiagnostics",
		Params: rr,
	}

	err := sendRPCResponse(w, rpcResponse)
	if err != nil {
		Log.Error(err)
	}
}

// Initialize gives a response with the server capabilities and info.
func Initialize(_ interface{}) (interface{}, error) {
	response := InitializeResult{
		Capabilities: ServerCapabilities{
			TextDocumentSync: Full,
		},
		ServerInfo: ServerInfo{
			Name:    "WhaleLintLSP",
			Version: "0.0.1",
		},
	}

	return response, nil
}

// Initialized is a handler for client's initialized notification
//
// As it does not contain an id, no response is expected.
func Initialized(_ []byte) (interface{}, error) {
	return nil, nil
}

func Shutdown(_ interface{}) (interface{}, error) {
	shutdownChannel <- true

	return nil, nil
}

func sendRPCResponse(w *bufio.Writer, rpcResponse interface{}) error {
	responseJSON, err := json.Marshal(rpcResponse)
	if err != nil {
		return fmt.Errorf("failed to marshal rpcResponse: %w", err)
	}

	contentLength := len(responseJSON)
	header := "Content-Length: " + strconv.Itoa(contentLength) + "\r\n\r\n"

	rawResponse := append([]byte(header), responseJSON...)

	Log.Debug("Send response to Client: ", string(rawResponse))

	_, err = w.Write(rawResponse)
	if err != nil {
		return fmt.Errorf("failed to send JSONRPC response: %w", err)
	}

	err = w.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush JSONRPC connection bufio.Writer: %w", err)
	}

	return nil
}

func HandleRequest(w *bufio.Writer, requestBytes []byte) error {
	Log.Debug("Client raw request: ", string(requestBytes))

	request := RPCRequest{}

	err := json.Unmarshal(requestBytes, &request)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSONRPC request : %w", err)
	}

	if request.ID == nil {
		// handle notifications
		// they do no require a response, but may trigger push from server
		Log.Debug("Received Notification from Client: ", request.Method)

		handler, ok := notificationHandlerMap[request.Method]
		if !ok {
			// unsupported call, that we do not handle at the moment.
			Log.Debug("Unsupported notification method:", request.Method)

			return nil
		}

		r, errNotification := handler(requestBytes)

		// publish r
		if rr, ok := r.(TextDocumentURIandStageList); ok {
			PublishDiagnostics(rr, w)
		}

		return errNotification
	}

	Log.Debug("Received Request from Client: ", request.Method)

	handler, ok := MethodMap[request.Method]
	if !ok {
		// unsupported call, that we do not handle at the moment.
		Log.Debug("Unsupported request method:", request.Method)

		return nil
	}

	response, errMethod := handler(request.Params)
	if errMethod != nil {
		Log.Error(errMethod)
	}

	rpcResponse := &RPCResponse{
		ID:     request.ID,
		Result: response,
		Err:    errMethod,
	}

	return sendRPCResponse(w, rpcResponse)
}

func HandleConnection(connection net.Conn, errC chan error) {
	log.Printf("Serving %s\n", connection.RemoteAddr().String())

	// buff := make([]byte, 5000)
	c := bufio.NewReader(connection)
	w := bufio.NewWriter(connection)

	// The base protocol consists of a header and a content part (comparable to HTTP).
	// The header and content part are separated by a ‘\r\n’.
	//
	// Content-Length: ..... \r\n
	// \r\n
	// [content part]

	for {
		// Content-length: .....
		// read n byte which contains the message length
		// read till /r/n
		contentLengthHeaderBytes, err := c.ReadBytes('\n')
		if err != nil {
			errC <- fmt.Errorf("failed to read JSONRPC request header bytes: %w", err)
		}

		contentLengthHeaderStr := string(contentLengthHeaderBytes[:len(contentLengthHeaderBytes)-2])

		// Skip the next "\r\n".
		_, err = c.Discard(len("\r\n"))
		if err != nil {
			errC <- fmt.Errorf("failed to discard 2 bytes (\r\n) from JSONRPC request : %w", err)
		}

		// Parse content length.
		contentLengthStr := strings.TrimPrefix(contentLengthHeaderStr, "Content-Length: ")
		contentLength, errAtoi := strconv.Atoi(contentLengthStr)
		if errAtoi != nil { // nolint:wsl
			errC <- fmt.Errorf("failed to parse content-length from string to int: %w", err)
		}

		requestContentBuff := make([]byte, contentLength)
		// Read the content.
		_, err = io.ReadFull(c, requestContentBuff)
		if err != nil {
			errC <- fmt.Errorf("failed to read JSONRPC content: %w", err)
		}

		err = HandleRequest(w, requestContentBuff)
		if err != nil {
			errC <- fmt.Errorf("failed to handle JSONRPC request: %w", err)
		}
	}
}

var shutdownChannel = make(chan bool) // nolint:gochecknoglobals

func Serve(port int) error {
	host := "0.0.0.0"

	// nolint:gofmt,gofumpt,goimports
	MethodMap = MethodMapType{
		"initialize": Initialize,
		"shutdown"  : Shutdown,
	}

	// nolint:gofmt,gofumpt,goimports
	notificationHandlerMap = NotificationHandlerMap{
		"initialized"           : Initialized,
		"textDocument/didOpen"  : OnTextOpen,
		"textDocument/didClose" : Yay,
		"textDocument/didChange": onTextDocumentDidChange,
		"textDocument/didSave"  : onTextDocumentDidSave,
	}

	serviceAddress := host + ":" + strconv.Itoa(port)

	tcpAddress, errResolveTCP := net.ResolveTCPAddr("tcp", serviceAddress)
	if errResolveTCP != nil {
		return fmt.Errorf("LSP | cannot resolve TCP address: %w", errResolveTCP)
	}

	listener, errListenTCP := net.ListenTCP("tcp", tcpAddress)
	if errListenTCP != nil {
		return fmt.Errorf("LSP | cannot listen on TCP address: %w", errListenTCP)
	}

	defer listener.Close()

	Log.Println("Listening on", tcpAddress)
	// c1 := make(chan string)

	errorChannel := make(chan error)
	shutdown := false

	// Serve
	for !shutdown {
		connection, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("LSP | cannot accept connection %w", err)
		}

		go HandleConnection(connection, errorChannel)

		select {
		case err := <-errorChannel:
			Log.Error(err)

			break
		case shutdown = <-shutdownChannel:
			continue
		default:
			continue
		}
	}

	return nil
}
