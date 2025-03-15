package protocol

import (
	"kro-extenstion/server/diagnostics"
	"kro-extenstion/server/validator"

	"github.com/tliron/commonlog"
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

	// Validate document when it's opened
	validateDocument(context, params.TextDocument.URI, params.TextDocument.Text)

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
	log := commonlog.GetLogger("validator")
	log.Infof("Validating document: %s", uri)
	log.Infof("Document content: %s", content)
	diagnosticManager.ClearDiagnostics(uri)

	// Parse YAML
	parser := validator.NewYAMLParser(content)
	data, err := parser.Parse()
	if err != nil {
		log.Errorf("Error parsing YAML: %s", err)
		diagnostic := CreateDiagnostic(
			err.Error(),
			protocol.DiagnosticSeverityError,
			CreateErrorRange(0, 0, 0),
		)
		diagnosticManager.AddDiagnostic(uri, diagnostic)
		// Make sure to publish diagnostics even on error
		log.Infof("Publishing diagnostics for parse error: %v", diagnostic)
		ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics,
			protocol.PublishDiagnosticsParams{
				URI:         uri,
				Diagnostics: diagnosticManager.GetDiagnostics(uri),
			})
		return
	}
	log.Infof("YAML parsed successfully: %v", data)
	// Validate schema
	errors := validator.ValidateResourceGraph(data)
	log.Infof("Validation errors: %v", errors)
	for _, err := range errors {
		diagnostic := CreateDiagnostic(
			err.Error(),
			protocol.DiagnosticSeverityError,
			CreateErrorRange(0, 0, 0),
		)
		diagnosticManager.AddDiagnostic(uri, diagnostic)
	}

	// Publish diagnostics using the passed context
	log.Infof("Publishing diagnostics: %v", diagnosticManager.GetDiagnostics(uri))
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics,
		protocol.PublishDiagnosticsParams{
			URI:         uri,
			Diagnostics: diagnosticManager.GetDiagnostics(uri),
		})
}
