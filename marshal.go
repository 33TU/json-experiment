package jsonexperiment

import (
	"sync"

	"github.com/33TU/json-experiment/internal"
)

type pooledBuffer struct {
	b []byte
}

var bytesPool = sync.Pool{
	New: func() any {
		return &pooledBuffer{
			b: make([]byte, 0, 1024),
		}
	},
}

// MarshalAppend appends the JSON encoding of value to dst.
func MarshalAppend(dst []byte, value any) ([]byte, error) {
	return MarshalAppendWithFlags(dst, value, 0)
}

// MarshalAppendWithOptions appends the JSON encoding of value to dst using opts.
func MarshalAppendWithOptions(dst []byte, value any, opts ...MarshalOptions) ([]byte, error) {
	flags := 0
	for _, opt := range opts {
		flags |= int(opt.Flags())
	}
	return MarshalAppendWithFlags(dst, value, MarshalFlags(flags))
}

// MarshalAppendWithFlags appends the JSON encoding of value to dst using flags.
func MarshalAppendWithFlags(dst []byte, value any, flags MarshalFlags) ([]byte, error) {
	start := len(dst)

	buf, err := marshalInterface(dst, value, flags)
	if err != nil || flags&MarshalFlagValidateString == 0 {
		return buf, err
	}

	return internal.ReplaceInvalidUTF8(buf, start), nil
}

// Marshal returns the JSON encoding of value.
func Marshal(value any) ([]byte, error) {
	return MarshalWithFlags(value, 0)
}

// MarshalWithOptions returns the JSON encoding of value using opts.
func MarshalWithOptions(value any, opts ...MarshalOptions) ([]byte, error) {
	flags := 0
	for _, opt := range opts {
		flags |= int(opt.Flags())
	}
	return MarshalWithFlags(value, MarshalFlags(flags))
}

// MarshalWithFlags returns the JSON encoding of value using flags.
func MarshalWithFlags(value any, flags MarshalFlags) ([]byte, error) {
	pb := bytesPool.Get().(*pooledBuffer)

	buf, err := marshalInterface(pb.b[:0], value, flags)
	if err == nil && flags&MarshalFlagValidateString != 0 {
		buf = internal.ReplaceInvalidUTF8(buf, 0)
	}
	out := append([]byte(nil), buf...)

	pb.b = buf[:0]
	bytesPool.Put(pb)

	return out, err
}
