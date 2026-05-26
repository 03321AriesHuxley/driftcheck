package drift

// DriftKind describes the type of drift detected.
type DriftKind string

const (
	DriftKindEnv     DriftKind = "env"
	DriftKindImage   DriftKind = "image"
	DriftKindPort    DriftKind = "port"
	DriftKindVolume  DriftKind = "volume"
	DriftKindCommand DriftKind = "command"
)

// DriftItem represents a single detected drift between a running container
// and its Compose service definition.
type DriftItem struct {
	Kind     DriftKind
	Field    string
	Expected string
	Actual   string
}

// Result holds the drift check outcome for a single service.
type Result struct {
	ServiceName   string
	ContainerID   string
	ContainerName string
	Drifts        []DriftItem
}

// HasDrift returns true when at least one drift item was recorded.
func (r *Result) HasDrift() bool {
	return len(r.Drifts) > 0
}

// AddDrift appends a new DriftItem to the result.
func (r *Result) AddDrift(kind DriftKind, field, expected, actual string) {
	r.Drifts = append(r.Drifts, DriftItem{
		Kind:     kind,
		Field:    field,
		Expected: expected,
		Actual:   actual,
	})
}
