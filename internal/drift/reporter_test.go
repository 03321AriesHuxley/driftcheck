package drift

import (
	"bytes"
	"strings"
	"testing"
)

func TestReporter_AllClean(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	results := []DriftResult{
		{ServiceName: "web", Clean: true},
		{ServiceName: "db", Clean: true},
	}
	r.Print(results)
	out := buf.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected clean message, got: %s", out)
	}
}

func TestReporter_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	results := []DriftResult{
		{ServiceName: "web", Clean: false, Drifts: []string{"image mismatch", "env PORT missing"}},
		{ServiceName: "db", Clean: true},
	}
	r.Print(results)
	out := buf.String()
	if !strings.Contains(out, "[web] ✗ drift detected") {
		t.Errorf("expected drift header for web, got: %s", out)
	}
	if !strings.Contains(out, "image mismatch") {
		t.Errorf("expected 'image mismatch' in output, got: %s", out)
	}
	if !strings.Contains(out, "[db] ✓ clean") {
		t.Errorf("expected db to be clean, got: %s", out)
	}
}

func TestReporter_Empty(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	r.Print([]DriftResult{})
	out := buf.String()
	if !strings.Contains(out, "No services checked") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestReporter_Summary(t *testing.T) {
	var buf bytes.Buffer
	r := NewReporter(&buf)
	results := []DriftResult{
		{ServiceName: "web", Clean: false, Drifts: []string{"image mismatch"}},
		{ServiceName: "db", Clean: true},
		{ServiceName: "cache", Clean: true},
	}
	summary := r.Summary(results)
	if !strings.Contains(summary, "3 service(s) checked") {
		t.Errorf("expected 3 services in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "1 drifted") {
		t.Errorf("expected 1 drifted in summary, got: %s", summary)
	}
	if !strings.Contains(summary, "2 clean") {
		t.Errorf("expected 2 clean in summary, got: %s", summary)
	}
}
