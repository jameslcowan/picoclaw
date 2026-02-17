package mcp

import (
	"context"
	"testing"
)

type fakeClient struct {
	tools      []RemoteTool
	callResult CallResult
	callErr    error

	lastToolName string
	lastArgs     map[string]any
}

func (f *fakeClient) Start(_ context.Context) error { return nil }
func (f *fakeClient) ListTools(_ context.Context) ([]RemoteTool, error) {
	return f.tools, nil
}
func (f *fakeClient) CallTool(_ context.Context, toolName string, arguments map[string]any) (CallResult, error) {
	f.lastToolName = toolName
	f.lastArgs = arguments
	if f.callErr != nil {
		return CallResult{}, f.callErr
	}
	return f.callResult, nil
}
func (f *fakeClient) Close() error { return nil }

func TestManager_DiscoverTools_FilterAndCall(t *testing.T) {
	serverCfg := map[string]ServerConfig{
		"Local Dev": {
			Command:      "fake",
			IncludeTools: []string{"alpha", "beta"},
			ExcludeTools: []string{"beta"},
		},
	}
	manager := NewManager(serverCfg)

	client := &fakeClient{
		tools: []RemoteTool{
			{Name: "alpha", Description: "tool alpha"},
			{Name: "beta", Description: "tool beta"},
			{Name: "gamma", Description: "tool gamma"},
		},
		callResult: CallResult{Content: "ok"},
	}
	manager.newClient = func(_ ServerConfig) Client {
		return client
	}

	tools, err := manager.DiscoverTools(context.Background())
	if err != nil {
		t.Fatalf("DiscoverTools() error = %v", err)
	}
	if len(tools) != 1 {
		t.Fatalf("DiscoverTools() returned %d tools, want 1", len(tools))
	}

	tool := tools[0]
	if tool.ToolName != "alpha" {
		t.Fatalf("discovered tool = %q, want alpha", tool.ToolName)
	}

	result, err := manager.CallTool(context.Background(), tool.QualifiedName, map[string]any{"x": 1})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}
	if result.Content != "ok" {
		t.Fatalf("CallTool() content = %q, want ok", result.Content)
	}
	if client.lastToolName != "alpha" {
		t.Fatalf("called MCP tool = %q, want alpha", client.lastToolName)
	}
}

func TestManager_NormalizeEmptySchema(t *testing.T) {
	serverCfg := map[string]ServerConfig{
		"srv": {Command: "fake"},
	}
	manager := NewManager(serverCfg)
	manager.newClient = func(_ ServerConfig) Client {
		return &fakeClient{
			tools: []RemoteTool{{Name: "empty_schema", InputSchema: nil}},
		}
	}

	tools, err := manager.DiscoverTools(context.Background())
	if err != nil {
		t.Fatalf("DiscoverTools() error = %v", err)
	}
	if len(tools) != 1 {
		t.Fatalf("DiscoverTools() returned %d tools, want 1", len(tools))
	}

	parameters := tools[0].Parameters
	if parameters["type"] != "object" {
		t.Fatalf("normalized schema type = %v, want object", parameters["type"])
	}
}
