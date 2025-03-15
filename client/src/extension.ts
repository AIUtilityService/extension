import * as path from "path";
import { workspace, ExtensionContext, window } from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  // Server executable path
  const serverPath = context.asAbsolutePath(
    path.join("..", "server", "my-language-server") // Go up one directory to find server folder
  );

  // Server options
  const serverOptions: ServerOptions = {
    run: {
      command: serverPath,
      transport: TransportKind.stdio,
    },
    debug: {
      command: serverPath,
      transport: TransportKind.stdio,
    },
  };

  // Client options
  const clientOptions: LanguageClientOptions = {
    documentSelector: [
      { scheme: "file", language: "mylanguage" },
      { scheme: "file", language: "yaml" }, // Add this for YAML support
    ],
    synchronize: {
      fileEvents: workspace.createFileSystemWatcher("**/*.{yaml,yml}"), // Update file watcher
    },
    outputChannel: window.createOutputChannel("My Language Server"),
  };
  // Create and start client
  client = new LanguageClient(
    "myLanguageServer",
    "My Language Server",
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
