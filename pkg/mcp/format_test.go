package mcp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestFormatCallPayload_TextAndStructured(t *testing.T) {
	raw := json.RawMessage(`{
		"content":[{"type":"text","text":"hello"}],
		"structuredContent":{"ok":true}
	}`)

	result, err := formatCallPayload(raw, 4096)
	if err != nil {
		t.Fatalf("formatCallPayload() error = %v", err)
	}
	if result.IsError {
		t.Fatalf("expected IsError=false")
	}
	if !strings.Contains(result.Content, "hello") {
		t.Fatalf("expected content to contain text block, got %q", result.Content)
	}
	if !strings.Contains(result.Content, `"ok":true`) {
		t.Fatalf("expected content to contain structured content, got %q", result.Content)
	}
}

func TestFormatCallPayload_Truncates(t *testing.T) {
	raw := json.RawMessage(`{"content":[{"type":"text","text":"abcdefghijklmnopqrstuvwxyz"}]}`)

	result, err := formatCallPayload(raw, 12)
	if err != nil {
		t.Fatalf("formatCallPayload() error = %v", err)
	}
	if len(result.Content) != 12 {
		t.Fatalf("expected truncated length 12, got %d", len(result.Content))
	}
}

func TestFormatCallPayload_RespectsIsError(t *testing.T) {
	raw := json.RawMessage(`{"content":[{"type":"text","text":"failed"}],"isError":true}`)

	result, err := formatCallPayload(raw, 4096)
	if err != nil {
		t.Fatalf("formatCallPayload() error = %v", err)
	}
	if !result.IsError {
		t.Fatalf("expected IsError=true")
	}
}
