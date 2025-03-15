package protocol

import (
	"kro-extenstion/server/diagnostics"
	"kro-extenstion/server/validator"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var documentManager = NewDocumentManager()
var diagnosticManager = diagnostics.NewDiagnosticManager()

func TextDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	documentManager.AddDocument(
		params.TextDocument.URI,
		params.TextDocument.Text,
		params.TextDocument.Version,
	)
	return nil
}

func TextDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) == 0 {
		return nil
	}

	textChange, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent)
	if !ok {
		return nil
	}

	// Update document
	documentManager.UpdateDocument(params.TextDocument.URI, textChange.Text, params.TextDocument.Version)

	// Pass the context to validateDocument
	validateDocument(context, params.TextDocument.URI, textChange.Text)

	return nil
}

func TextDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	documentManager.RemoveDocument(params.TextDocument.URI)
	return nil
}

// Update the function signature to include context parameter
func validateDocument(ctx *glsp.Context, uri string, content string) {
	diagnosticManager.ClearDiagnostics(uri)

	// Parse YAML
	parser := validator.NewYAMLParser(content)
	data, err := parser.Parse()
	if err != nil {
		diagnostic := CreateDiagnostic(
			err.Error(),
			protocol.DiagnosticSeverityError,
			CreateErrorRange(0, 0, 0),
		)
		diagnosticManager.AddDiagnostic(uri, diagnostic)
		return
	}

	// Validate schema
	errors := validator.ValidateResourceGraph(data)
	for _, err := range errors {
		diagnostic := CreateDiagnostic(
			err.Error(),
			protocol.DiagnosticSeverityError,
			CreateErrorRange(0, 0, 0),
		)
		diagnosticManager.AddDiagnostic(uri, diagnostic)
	}

	// Publish diagnostics using the passed context
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics,
		protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnosticManager.GetDiagnostics(uri),
		})
}
