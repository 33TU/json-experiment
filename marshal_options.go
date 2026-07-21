package jsonexperiment

// MarshalOptions configures JSON marshaling behavior.
type MarshalOptions struct {
	// EscapeHTML escapes <, >, and & as Unicode escape sequences inside strings.
	EscapeHTML bool
}

type marshalFlags uint32

const (
	marshalFlagEscapeHTML marshalFlags = 1 << iota
)

func (o MarshalOptions) flags() marshalFlags {
	var flags marshalFlags
	if o.EscapeHTML {
		flags |= marshalFlagEscapeHTML
	}
	return flags
}
