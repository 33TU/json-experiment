package jsonexperiment

// MarshalOptions configures JSON marshaling behavior.
type MarshalOptions struct {
	// EscapeHTML escapes <, >, and & as Unicode escape sequences inside strings.
	EscapeHTML bool

	// ValidateString replaces invalid UTF-8 bytes with the Unicode replacement character.
	ValidateString bool

	// SortMapKeys sorts map keys by their encoded JSON object name.
	SortMapKeys bool
}

// MarshalFlags represents flags that control the behavior of the marshaling process.
type MarshalFlags uint32

const (
	MarshalFlagEscapeHTML MarshalFlags = 1 << iota
	MarshalFlagValidateString
	MarshalFlagSortMapKeys
)

func (o MarshalOptions) Flags() MarshalFlags {
	var flags MarshalFlags
	if o.EscapeHTML {
		flags |= MarshalFlagEscapeHTML
	}
	if o.ValidateString {
		flags |= MarshalFlagValidateString
	}
	if o.SortMapKeys {
		flags |= MarshalFlagSortMapKeys
	}
	return flags
}
