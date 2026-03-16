package api_test

import (
	"path/filepath"
	"testing"

	"github.com/organic-programming/grace-op/api"
	opv1 "github.com/organic-programming/grace-op/gen/go/op/v1"
)

func TestCheckReturnsLifecycleReport(t *testing.T) {
	root := t.TempDir()
	dir := writeProtoHolon(t, root)

	resp, err := api.Check(&opv1.LifecycleRequest{Target: dir})
	if err != nil {
		t.Fatalf("Check error = %v", err)
	}
	if got := resp.GetReport().GetOperation(); got != "check" {
		t.Fatalf("operation = %q, want %q", got, "check")
	}
	if got := filepath.Base(resp.GetReport().GetManifest()); got != "holon.proto" {
		t.Fatalf("manifest basename = %q, want %q", got, "holon.proto")
	}
}
