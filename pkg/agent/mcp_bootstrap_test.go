package agent

import (
	"testing"
	"time"

	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/mcp"
)

func TestCalculateMCPDiscoveryTimeout_UsesMaxInitWithGrace(t *testing.T) {
	serverConfigs := map[string]struct {
		initSeconds int
	}{
		"fast": {initSeconds: 5},
		"slow": {initSeconds: 60},
	}

	cfg := config.MCPToolsConfig{
		Enabled: true,
		Servers: map[string]config.MCPServerConfig{
			"fast": {
				Enabled:            true,
				Command:            "fast",
				InitTimeoutSeconds: serverConfigs["fast"].initSeconds,
			},
			"slow": {
				Enabled:            true,
				Command:            "slow",
				InitTimeoutSeconds: serverConfigs["slow"].initSeconds,
			},
		},
	}

	mcpConfigs := buildMCPServerConfigs(cfg)
	timeout := calculateMCPDiscoveryTimeout(mcpConfigs)

	want := 65 * time.Second
	if timeout != want {
		t.Fatalf("calculateMCPDiscoveryTimeout() = %v, want %v", timeout, want)
	}
}

func TestBuildMCPServerConfigs_SkipsDisabledServers(t *testing.T) {
	cfg := config.MCPToolsConfig{
		Enabled: true,
		Servers: map[string]config.MCPServerConfig{
			"context7": {
				Enabled:  true,
				Command:  "context7-mcp",
				Protocol: "jsonl",
			},
			"disabled": {
				Enabled: false,
				Command: "ignored",
			},
		},
	}

	mcpConfigs := buildMCPServerConfigs(cfg)
	if len(mcpConfigs) != 1 {
		t.Fatalf("buildMCPServerConfigs() count = %d, want 1", len(mcpConfigs))
	}

	context7, ok := mcpConfigs["context7"]
	if !ok {
		t.Fatalf("context7 not found in buildMCPServerConfigs output")
	}
	if context7.Protocol != "jsonl" {
		t.Fatalf("context7 protocol = %q, want jsonl", context7.Protocol)
	}
}

func TestInferMCPProtocol_Context7DefaultsToJSONL(t *testing.T) {
	got := inferMCPProtocol("", "context7-mcp")
	if got != mcp.ProtocolJSONLines {
		t.Fatalf("inferMCPProtocol() = %q, want %s", got, mcp.ProtocolJSONLines)
	}
}
