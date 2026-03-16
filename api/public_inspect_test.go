package api_test

import (
	"testing"

	"github.com/organic-programming/grace-op/api"
	opv1 "github.com/organic-programming/grace-op/gen/go/op/v1"
)

func TestInspectLoadsLocalProtoHolon(t *testing.T) {
	root := t.TempDir()
	dir := writeProtoHolon(t, root)

	resp, err := api.Inspect(&opv1.InspectRequest{Target: dir})
	if err != nil {
		t.Fatalf("Inspect error = %v", err)
	}
	if len(resp.GetDocument().GetServices()) != 1 {
		t.Fatalf("services = %d, want 1", len(resp.GetDocument().GetServices()))
	}
	service := resp.GetDocument().GetServices()[0]
	if service.GetName() != "alpha.v1.AlphaService" {
		t.Fatalf("service name = %q, want %q", service.GetName(), "alpha.v1.AlphaService")
	}
	if len(service.GetMethods()) != 1 || service.GetMethods()[0].GetName() != "Ping" {
		t.Fatalf("methods = %#v, want Ping", service.GetMethods())
	}
}
