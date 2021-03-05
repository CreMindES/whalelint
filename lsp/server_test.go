package lsp_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"

	LSP "github.com/cremindes/whalelint/lsp"
)

func TestOnTextOpen(t *testing.T) {
	t.Parallel()

	str := "FROM golang:1.16"

	reader := strings.NewReader(str)

	dockerfile, parseErr := parser.Parse(reader)
	assert.Nil(t, parseErr)

	stageList, _, paerseStageErr := instructions.Parse(dockerfile.AST)
	assert.Nil(t, paerseStageErr)

	expected := LSP.TextDocumentURIandStageList{
		StageList: stageList,
		URI:       "mockURI",
	}

	type TextDocumentWrapper struct {
		TextDocument LSP.TextDocumentItem `json:"textDocument"`
	}

	testDocParam := struct {
		JSONrpcVersion string              `json:"jsonrpc"`      // "jsonrpc": "2.0",
		ID             interface{}         `json:"id,omitempty"` // "id": 1,
		Method         string              `json:"method"`       // "method": "textDocument/didOpen",
		Params         TextDocumentWrapper `json:"params"`
	}{
		"2.0",
		1,
		"textDocument/didOpen",
		TextDocumentWrapper{LSP.TextDocumentItem{ // nolint:exhaustivestruct
			URI:        "mockURI",
			LanguageID: "dockerfile",
			Version:    0,
			Text:       str,
		}},
	}

	json, jsonErr := json.Marshal(testDocParam)
	assert.Nil(t, jsonErr)

	result, err := LSP.OnTextOpen(json)

	assert.Equal(t, expected, result)
	assert.Nil(t, err)
}
