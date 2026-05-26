package drift

import (
	"testing"
)

func TestResult_HasDrift_False(t *testing.T) {
	r := &Result{
		ServiceName:   "web",
		ContainerID:   "abc123",
		ContainerName: "myapp_web_1",
	}
	if r.HasDrift() {
		t.Error("expected no drift on empty result")
	}
}

func TestResult_HasDrift_True(t *testing.T) {
	r := &Result{
		ServiceName: "web",
	}
	r.AddDrift(DriftKindEnv, "PORT", "8080", "9090")
	if !r.HasDrift() {
		t.Error("expected drift to be detected")
	}
}

func TestResult_AddDrift_MultipleItems(t *testing.T) {
	r := &Result{ServiceName: "db"}
	r.AddDrift(DriftKindImage, "image", "postgres:14", "postgres:13")
	r.AddDrift(DriftKindEnv, "POSTGRES_DB", "mydb", "otherdb")

	if len(r.Drifts) != 2 {
		t.Fatalf("expected 2 drift items, got %d", len(r.Drifts))
	}

	if r.Drifts[0].Kind != DriftKindImage {
		t.Errorf("expected first drift kind %q, got %q", DriftKindImage, r.Drifts[0].Kind)
	}
	if r.Drifts[1].Field != "POSTGRES_DB" {
		t.Errorf("expected field POSTGRES_DB, got %q", r.Drifts[1].Field)
	}
}

func TestResult_AddDrift_Fields(t *testing.T) {
	r := &Result{ServiceName: "cache"}
	r.AddDrift(DriftKindPort, "6379/tcp", "6379:6379", "")

	item := r.Drifts[0]
	if item.Expected != "6379:6379" {
		t.Errorf("unexpected Expected value: %q", item.Expected)
	}
	if item.Actual != "" {
		t.Errorf("unexpected Actual value: %q", item.Actual)
	}
}
