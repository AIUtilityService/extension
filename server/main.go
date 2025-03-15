package main

import (
	"kro-extenstion/protocol"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	lspProtocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const lsName = "kro-language-server"

var (
	version string = "0.0.1"
	handler lspProtocol.Handler
)

func main() {
	// Configure logging
	commonlog.Configure(1, nil)
	log := commonlog.GetLogger("server")

	handler = lspProtocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,
		TextDocumentDidOpen: func(context *glsp.Context, params *lspProtocol.DidOpenTextDocumentParams) error {
			log.Infof("Document opened: %s", params.TextDocument.URI)
			return protocol.TextDocumentDidOpen(context, params)
		},
		TextDocumentDidChange: func(context *glsp.Context, params *lspProtocol.DidChangeTextDocumentParams) error {
			log.Infof("Document changed: %s", params.TextDocument.URI)
			return protocol.TextDocumentDidChange(context, params)
		},
		TextDocumentDidClose: protocol.TextDocumentDidClose,
	}

	server := server.NewServer(&handler, lsName, false)
	server.RunStdio()
}

func initialize(context *glsp.Context, params *lspProtocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	openClose := true
	changeValue := lspProtocol.TextDocumentSyncKindFull

	// Set specific capabilities with manual pointer values
	capabilities.TextDocumentSync = &lspProtocol.TextDocumentSyncOptions{
		OpenClose: &openClose,
		Change:    &changeValue,
	}

	return lspProtocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &lspProtocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *lspProtocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	return nil
}
