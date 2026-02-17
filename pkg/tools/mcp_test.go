package tools

import (
	"testing"

	"github.com/sipeed/picoclaw/pkg/mcp"
)

func TestRegisterKnownMCPTools_RegistersAllTools(t *testing.T) {
	registry := NewToolRegistry()
	manager := &mcp.Manager{}
	discovered := []mcp.RegisteredTool{
		{
			QualifiedName: "mcp_context7__resolve_library_id",
			ServerName:    "context7",
			ToolName:      "resolve-library-id",
			Description:   "Resolve library ID",
			Parameters: map[string]any{
				"type": "object",
			},
		},
		{
			QualifiedName: "mcp_context7__query_docs",
			ServerName:    "context7",
			ToolName:      "query-docs",
			Description:   "Query docs",
			Parameters: map[string]any{
				"type": "object",
			},
		},
	}

	count := RegisterKnownMCPTools(registry, manager, discovered)
	if count != 2 {
		t.Fatalf("RegisterKnownMCPTools count = %d, want 2", count)
	}

	if _, ok := registry.Get("mcp_context7__resolve_library_id"); !ok {
		t.Fatalf("expected mcp_context7__resolve_library_id to be registered")
	}
	if _, ok := registry.Get("mcp_context7__query_docs"); !ok {
		t.Fatalf("expected mcp_context7__query_docs to be registered")
	}
}
