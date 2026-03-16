package api_test

import (
	"strings"
	"testing"

	"github.com/organic-programming/grace-op/api"
	opv1 "github.com/organic-programming/grace-op/gen/go/op/v1"
)

func TestRunSequenceDryRunRendersParams(t *testing.T) {
	root := t.TempDir()
	dir := writeProtoHolon(t, root)

	resp, err := api.RunSequence(&opv1.RunSequenceRequest{
		Holon:    dir,
		Sequence: "greet",
		DryRun:   true,
		Params:   map[string]string{"name": "Maria"},
	})
	if err != nil {
		t.Fatalf("RunSequence error = %v", err)
	}
	if got := len(resp.GetResult().GetSteps()); got != 1 {
		t.Fatalf("steps = %d, want 1", got)
	}
	if got := resp.GetResult().GetSteps()[0].GetCommand(); !strings.Contains(got, "Maria") {
		t.Fatalf("command = %q, want rendered param", got)
	}
}
