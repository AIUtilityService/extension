package protocol

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var documentManager = NewDocumentManager()

func TextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	documentManager.AddDocument(
		params.TextDocument.URI,
		params.TextDocument.Text,
		params.TextDocument.Version,
	)
	return nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	// For full sync, we just need the last change
	if len(params.ContentChanges) > 0 {
		documentManager.UpdateDocument(
			params.TextDocument.URI,
			params.ContentChanges[len(params.ContentChanges)-1].Text,
			params.TextDocument.Version,
		)
	}
	return nil
}

func TextDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	documentManager.RemoveDocument(params.TextDocument.URI)
	return nil
}
