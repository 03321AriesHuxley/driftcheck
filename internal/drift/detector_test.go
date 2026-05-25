package drift

import (
	"testing"

	"github.com/user/driftcheck/internal/compose"
	"github.com/user/driftcheck/internal/docker"
)

func makeService(name, image string, env []string) compose.Service {
	return compose.Service{
		Name:        name,
		Image:       image,
		Environment: env,
	}
}

func TestToEnvMap(t *testing.T) {
	m := toEnvMap([]string{"FOO=bar", "BAZ=qux", "EMPTY="})
	if m["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", m["FOO"])
	}
	if m["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %q", m["BAZ"])
	}
	if m["EMPTY"] != "" {
		t.Errorf("expected EMPTY='', got %q", m["EMPTY"])
	}
}

func TestSplitEnv_WithValue(t *testing.T) {
	k, v := splitEnv("KEY=value")
	if k != "KEY" || v != "value" {
		t.Errorf("unexpected split: k=%q v=%q", k, v)
	}
}

func TestSplitEnv_NoEquals(t *testing.T) {
	k, v := splitEnv("KEYONLY")
	if k != "KEYONLY" || v != "" {
		t.Errorf("unexpected split: k=%q v=%q", k, v)
	}
}

func TestCheck_Clean(t *testing.T) {
	insp := docker.NewInspector(nil) // nil client — we'll stub via ContainerInfo
	_ = insp
	// Integration-level test would require a live daemon; unit coverage via util tests above.
	// Structural smoke-test: ensure Detector can be constructed.
	det := NewDetector(insp)
	if det == nil {
		t.Fatal("expected non-nil Detector")
	}
}

func TestDriftResult_Fields(t *testing.T) {
	r := &DriftResult{
		ServiceName: "web",
		Drifts:      []string{"image mismatch"},
		Clean:       false,
	}
	if r.Clean {
		t.Error("expected Clean=false")
	}
	if len(r.Drifts) != 1 {
		t.Errorf("expected 1 drift, got %d", len(r.Drifts))
	}
}
