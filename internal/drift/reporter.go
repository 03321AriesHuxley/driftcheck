package drift

import (
	"fmt"
	"io"
	"strings"
)

// DriftResult holds the outcome of a drift check for a single service.
type DriftResult struct {
	ServiceName string
	Drifts      []string
	Clean       bool
}

// Reporter formats and writes drift results to an output stream.
type Reporter struct {
	w io.Writer
}

// NewReporter creates a Reporter that writes to w.
func NewReporter(w io.Writer) *Reporter {
	return &Reporter{w: w}
}

// Print writes a human-readable summary of the drift results.
func (r *Reporter) Print(results []DriftResult) {
	if len(results) == 0 {
		fmt.Fprintln(r.w, "No services checked.")
		return
	}

	allClean := true
	for _, res := range results {
		if !res.Clean {
			allClean = false
			break
		}
	}

	if allClean {
		fmt.Fprintln(r.w, "✓ No drift detected — all services match their Compose definitions.")
		return
	}

	for _, res := range results {
		if res.Clean {
			fmt.Fprintf(r.w, "[%s] ✓ clean\n", res.ServiceName)
			continue
		}
		fmt.Fprintf(r.w, "[%s] ✗ drift detected:\n", res.ServiceName)
		for _, d := range res.Drifts {
			fmt.Fprintf(r.w, "  - %s\n", d)
		}
	}
}

// Summary returns a single-line summary string.
func (r *Reporter) Summary(results []DriftResult) string {
	total := len(results)
	drifted := 0
	for _, res := range results {
		if !res.Clean {
			drifted++
		}
	}
	parts := []string{
		fmt.Sprintf("%d service(s) checked", total),
		fmt.Sprintf("%d drifted", drifted),
		fmt.Sprintf("%d clean", total-drifted),
	}
	return strings.Join(parts, ", ")
}
