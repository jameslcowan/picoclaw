package mcp

import "testing"

func TestNormalizeProtocol(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ProtocolMCPFrames},
		{input: "mcp", want: ProtocolMCPFrames},
		{input: "jsonl", want: ProtocolJSONLines},
		{input: "JSONL", want: ProtocolJSONLines},
		{input: "unknown", want: ProtocolMCPFrames},
	}

	for _, tt := range tests {
		got := normalizeProtocol(tt.input)
		if got != tt.want {
			t.Fatalf("normalizeProtocol(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
