package api_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/organic-programming/grace-op/api"
	opv1 "github.com/organic-programming/grace-op/gen/go/op/v1"
)

func TestToolsReturnsOpenAIPayload(t *testing.T) {
	root := t.TempDir()
	dir := writeProtoHolon(t, root)

	resp, err := api.Tools(&opv1.ToolsRequest{Target: dir, Format: "openai"})
	if err != nil {
		t.Fatalf("Tools error = %v", err)
	}
	if got := resp.GetFormat(); got != "openai" {
		t.Fatalf("format = %q, want %q", got, "openai")
	}
	if !strings.Contains(string(resp.GetPayload()), "Ping") {
		t.Fatalf("payload = %q, want to mention Ping", string(resp.GetPayload()))
	}
}

func TestEnvInitializesDirectoriesAndShell(t *testing.T) {
	root := t.TempDir()
	oppath := filepath.Join(root, ".op")
	opbin := filepath.Join(oppath, "bin")
	t.Setenv("OPPATH", oppath)
	t.Setenv("OPBIN", opbin)

	resp, err := api.Env(&opv1.EnvRequest{Init: true, Shell: true})
	if err != nil {
		t.Fatalf("Env error = %v", err)
	}
	if resp.GetOppath() != oppath {
		t.Fatalf("OPPATH = %q, want %q", resp.GetOppath(), oppath)
	}
	if resp.GetOpbin() != opbin {
		t.Fatalf("OPBIN = %q, want %q", resp.GetOpbin(), opbin)
	}
	if _, err := os.Stat(opbin); err != nil {
		t.Fatalf("OPBIN directory missing: %v", err)
	}
	if !strings.Contains(resp.GetShell(), "export OPPATH") {
		t.Fatalf("shell snippet = %q, want export OPPATH", resp.GetShell())
	}
}
