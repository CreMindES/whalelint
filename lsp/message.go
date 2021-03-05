package lsp

import "encoding/json"

// Copied from https://github.com/golang/tools/blob/master/internal/lsp/protocol/tsprotocol.go.
// TODO: generate this automatically

/*
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

type RPCRequest struct {
	JSONrpcVersion string                 `json:"jsonrpc"`      // "jsonrpc": "2.0",
	ID             interface{}            `json:"id,omitempty"` // "id": 1,
	Method         string                 `json:"method"`       // "method": "textDocument/didOpen",
	Params         map[string]interface{} `json:"params"`       // "params": { ... }
}

type RPCResponse struct {
	ID interface{} `json:"id"`
	//	JSONRpcV string      `json:"jsonrpc"` // no need to have/set this fields, as it's constant at the moment.
	Result interface{} `json:"result"`
	Err    interface{} `json:"error"`
}

type RPCNotification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

func (r *RPCResponse) MarshalJSON() ([]byte, error) {
	if r.Err == nil {
		res := struct {
			JSONRpcV string      `json:"jsonrpc"`
			ID       interface{} `json:"id"`
			Result   interface{} `json:"result"`
		}{
			JSONRpcV: "2.0",
			ID:       r.ID,
			Result:   r.Result,
		}

		return json.Marshal(res)
	}

	return json.Marshal(r)
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities = struct {
	/**
	 * Defines how text documents are synced. Is either a detailed structure defining each notification or
	 * for backwards compatibility the TextDocumentSyncKind number.
	 */
	TextDocumentSync interface{} /*TextDocumentSyncOptions | TextDocumentSyncKind*/ `json:"textDocumentSync,omitempty"`
}

/**
 * Defines the capabilities provided by a language
 * server.
 */
type InnerServerCapabilities struct {
	TextDocumentSync interface{} /*TextDocumentSyncOptions | TextDocumentSyncKind*/ `json:"textDocumentSync,omitempty"`
}

type TextDocumentSyncKind float64

type TextDocumentSyncOptions struct { // nolint:maligned
	/**
	 * Open and close notifications are sent to the server. If omitted open close notification should not
	 * be sent.
	 */
	OpenClose bool `json:"openClose,omitempty"`
	/**
	 * Change notifications are sent to the server. See TextDocumentSyncKind.None, TextDocumentSyncKind.Full
	 * and TextDocumentSyncKind.Incremental. If omitted it defaults to TextDocumentSyncKind.None.
	 */
	Change TextDocumentSyncKind `json:"change,omitempty"`
	/**
	 * If present will save notifications are sent to the server. If omitted the notification should not be
	 * sent.
	 */
	WillSave bool `json:"willSave,omitempty"`
	/**
	 * If present will save wait until requests are sent to the server. If omitted the request should not be
	 * sent.
	 */
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	/**
	 * If present save notifications are sent to the server. If omitted the notification should not be
	 * sent.
	 */
	Save SaveOptions/*boolean | SaveOptions*/ `json:"save,omitempty"`
}

/**
 * Save options.
 */
type SaveOptions struct {
	/**
	 * The client is supposed to include the content on save.
	 */
	IncludeText bool `json:"includeText,omitempty"`
}

const (
	None TextDocumentSyncKind = 0
	/**
	 * Documents are synced by always sending the full content
	 * of the document.
	 */

	Full TextDocumentSyncKind = 1
	/**
	 * Documents are synced by sending the full content on open.
	 * After that only incremental updates to the document are
	 * send.
	 */

	Incremental TextDocumentSyncKind = 2
)

/**
 * A tagging type for string properties that are actually document URIs.
 */
type DocumentURI string

/**
 * An item to transfer a text document from the client to the
 * server.
 */
type TextDocumentItem struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentURI `json:"URI"`
	/**
	 * The text document's language identifier
	 */
	LanguageID string `json:"languageId"`
	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version float64 `json:"version"`
	/**
	 * The content of the opened text document.
	 */
	Text string `json:"text"`
}

/**
 * The change text document notification's parameters.
 */
type DidChangeTextDocumentParams struct {
	/**
	 * The document that did change. The version number points
	 * to the version after all provided content changes have
	 * been applied.
	 */
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	/**
	 * The actual content changes. The content changes describe single state changes
	 * to the document. So if there are two content changes c1 (at array index 0) and
	 * c2 (at array index 1) for a document in state S then c1 moves the document from
	 * S to S' and c2 from S' to S''. So c1 is computed on the state S and c2 is computed
	 * on the state S'.
	 *
	 * To mirror the content of a document using change events use the following approach:
	 * - start with the same initial content
	 * - apply the 'textDocument/didChange' notifications in the order you receive them.
	 * - apply the `TextDocumentContentChangeEvent`s in a single notification in the order
	 *   you receive them.
	 */
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

/**
 * A text document identifier to denote a specific version of a text document.
 */
type VersionedTextDocumentIdentifier struct {
	/**
	 * The version number of this document.
	 */
	Version int32 `json:"version"`
	TextDocumentIdentifier
}

/**
 * An event describing a change to a text document. If range and rangeLength are omitted
 * the new text is considered to be the full content of the document.
 *
 * @deprecated Use the text document from the new vscode-languageserver-textdocument package.
 */
type TextDocumentContentChangeEvent = struct {
	/**
	 * The range of the document that changed.
	 */
	Range *Range `json:"range,omitempty"`
	/**
	 * The optional length of the range that got replaced.
	 *
	 * @deprecated use range instead.
	 */
	RangeLength uint32 `json:"rangeLength,omitempty"`
	/**
	 * The new text for the provided range.
	 */
	Text string `json:"text"`
}

/**
 * A literal to identify a text document in the client.
 */
type TextDocumentIdentifier struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentURI `json:"URI"`
}

/**
 * The publish diagnostic notification's parameters.
 */
type PublishDiagnosticsParams struct {
	/**
	 * The URI for which diagnostic information is reported.
	 */
	URI DocumentURI `json:"URI"`
	/**
	 * Optional the version number of the document the diagnostics are published for.
	 *
	 * @since 3.15.0
	 */
	Version float64 `json:"version,omitempty"`
	/**
	 * An array of diagnostic information items.
	 */
	Diagnostics []Diagnostic `json:"diagnostics"`
}

/**
 * The diagnostic's severity.
 */
type DiagnosticSeverity float64

type Diagnostic struct {
	/**
	 * The range at which the message applies
	 */
	Range Range `json:"range"`
	/**
	 * The diagnostic's severity. Can be omitted. If omitted it is up to the
	 * client to interpret diagnostics as error, warning, info or hint.
	 */
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	/**
	 * The diagnostic's code, which usually appear in the user interface.
	 */
	Code interface{}/*number | string*/ `json:"code,omitempty"`
	/**
	 * An optional property to describe the error code.
	 *
	 * @since 3.16.0 - proposed state
	 */
	CodeDescription *CodeDescription `json:"codeDescription,omitempty"`
	/**
	 * A human-readable string describing the source of this
	 * diagnostic, e.g. 'typescript' or 'super lint'. It usually
	 * appears in the user interface.
	 */
	Source string `json:"source,omitempty"`
	/**
	 * The diagnostic's message. It usually appears in the user interface
	 */
	Message string `json:"message"`
	/**
	 * Additional metadata about the diagnostic.
	 *
	 * @since 3.15.0
	 */
	Tags []DiagnosticTag `json:"tags,omitempty"`
	/**
	 * An array of related diagnostic information, e.g. when symbol-names within
	 * a scope collide all definitions can be marked via this property.
	 */
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	/**
	 * A data entry field that is preserved between a `textDocument/publishDiagnostics`
	 * notification and `textDocument/codeAction` request.
	 *
	 * @since 3.16.0 - proposed state
	 */
	Data interface{} `json:"data,omitempty"`
}

/**
 * The diagnostic tags.
 *
 * @since 3.15.0.
 */
type DiagnosticTag float64

type Range struct {
	/**
	 * The range's start position
	 */
	Start Position `json:"start"`
	/**
	 * The range's end position.
	 */
	End Position `json:"end"`
}

/**
 * Position in a text document expressed as zero-based line and character offset.
 * The offsets are based on a UTF-16 string representation. So a string of the form
 * `a𐐀b` the character offset of the character `a` is 0, the character offset of `𐐀`
 * is 1 and the character offset of b is 3 since `𐐀` is represented using two code
 * units in UTF-16.
 *
 * Positions are line end character agnostic. So you can not specify a position that
 * denotes `\r|\n` or `\n|` where `|` represents the character offset.
 */
type Position struct {
	/**
	 * Line position in a document (zero-based).
	 * If a line number is greater than the number of lines in a document, it defaults back to the number of lines in
	 * the document.
	 * If a line number is negative, it defaults to 0.
	 */
	Line float64 `json:"line"`
	/**
	 * Character offset on a line in a document (zero-based). Assuming that the line is
	 * represented as a string, the `character` value represents the gap between the
	 * `character` and `character + 1`.
	 *
	 * If the character value is greater than the line length it defaults back to the
	 * line length.
	 * If a line number is negative, it defaults to 0.
	 */
	Character float64 `json:"character"`
}

/**
 * Structure to capture a description for an error code.
 *
 * @since 3.16.0 - proposed state.
 */
type CodeDescription struct {
	/**
	 * An URI to open with more information about the diagnostic error.
	 */
	Href URI `json:"href"`
}

/**
 * A tagging type for string properties that are actually URIs
 *
 * @since 3.16.0 - proposed state.
 */
type URI = string

/**
 * Represents a related message and source code location for a diagnostic. This should be
 * used to point to code locations that cause or related to a diagnostics, e.g when duplicating
 * a symbol in a scope.
 */
type DiagnosticRelatedInformation struct {
	/**
	 * The location of this related diagnostic information.
	 */
	Location Location `json:"location"`
	/**
	 * The message of this related diagnostic information.
	 */
	Message string `json:"message"`
}

/**
 * Represents a location inside a resource, such as a line
 * inside a text file.
 */
type Location struct {
	URI   DocumentURI `json:"URI"`
	Range Range       `json:"range"`
}

const (
	SeverityError DiagnosticSeverity = 1
	/**
	 * Reports a warning.
	 */

	SeverityWarning DiagnosticSeverity = 2
	/**
	 * Reports an information.
	 */

	SeverityInformation DiagnosticSeverity = 3
	/**
	 * Reports a hint.
	 */

	SeverityHint DiagnosticSeverity = 4
)

const (
	Unnecessary DiagnosticTag = 1
	/**
	 * Deprecated or obsolete code.
	 *
	 * Clients are allowed to rendered diagnostics with this tag strike through.
	 */

	Deprecated DiagnosticTag = 2
	/**
	 * A textual occurrence.
	 */
)
