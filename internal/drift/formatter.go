package drift

import (
	"fmt"
	"io"
	"strings"
)

// FormatResults writes a human-readable drift report to w.
// It returns the number of services with drift and any write error.
func FormatResults(w io.Writer, results []Result) (int, error) {
	driftCount := 0

	for _, r := range results {
		if !r.HasDrift() {
			_, err := fmt.Fprintf(w, "[OK]    %s (%s)\n", r.ServiceName, r.ContainerName)
			if err != nil {
				return driftCount, err
			}
			continue
		}

		driftCount++
		_, err := fmt.Fprintf(w, "[DRIFT] %s (%s)\n", r.ServiceName, r.ContainerName)
		if err != nil {
			return driftCount, err
		}

		for _, d := range r.Drifts {
			expected := d.Expected
			if expected == "" {
				expected = "<none>"
			}
			actual := d.Actual
			if actual == "" {
				actual = "<none>"
			}
			_, err = fmt.Fprintf(w, "        [%s] %s: expected=%q actual=%q\n",
				strings.ToUpper(string(d.Kind)), d.Field, expected, actual)
			if err != nil {
				return driftCount, err
			}
		}
	}

	return driftCount, nil
}
