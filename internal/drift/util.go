package drift

import "strings"

// toEnvMap converts a slice of "KEY=VALUE" strings into a map.
func toEnvMap(env []string) map[string]string {
	m := make(map[string]string, len(env))
	for _, kv := range env {
		k, v := splitEnv(kv)
		m[k] = v
	}
	return m
}

// splitEnv splits a "KEY=VALUE" string into its key and value parts.
// If there is no "=" the value is returned as an empty string.
func splitEnv(kv string) (string, string) {
	parts := strings.SplitN(kv, "=", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
