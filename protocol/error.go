package protocol

import (
	lspProtocol "github.com/tliron/glsp/protocol_3_16"
)

func CreateDiagnostic(message string, severity lspProtocol.DiagnosticSeverity, rng lspProtocol.Range) lspProtocol.Diagnostic {
	source := "kro-language-server" // Create a variable for the source string
	severityValue := severity       // Create a variable for severity

	return lspProtocol.Diagnostic{
		Range:    rng,
		Severity: &severityValue, // Assign severity as a pointer
		Source:   &source,        // Assign source as a pointer
		Message:  message,
	}
}

func CreateErrorRange(line int, startChar int, endChar int) lspProtocol.Range {
	return lspProtocol.Range{
		Start: lspProtocol.Position{
			Line:      uint32(line),
			Character: uint32(startChar),
		},
		End: lspProtocol.Position{
			Line:      uint32(line),
			Character: uint32(endChar),
		},
	}
}
