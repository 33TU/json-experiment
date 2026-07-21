package jsonexperiment

import "sync"

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

func MarshalAppend(dst []byte, value any) ([]byte, error) {
	return marshalInterface(dst, value, 0)
}

// MarshalAppendWithOptions appends the JSON encoding of value to dst using options.
func MarshalAppendWithOptions(dst []byte, value any, options MarshalOptions) ([]byte, error) {
	return marshalInterface(dst, value, options.flags())
}

func Marshal(value any) ([]byte, error) {
	pb := bytesPool.Get().(*pooledBuffer)

	buf, err := marshalInterface(pb.b[:0], value, 0)
	out := append([]byte(nil), buf...)

	pb.b = buf[:0]
	bytesPool.Put(pb)

	return out, err
}

// MarshalWithOptions returns the JSON encoding of value using options.
func MarshalWithOptions(value any, options MarshalOptions) ([]byte, error) {
	pb := bytesPool.Get().(*pooledBuffer)

	buf, err := marshalInterface(pb.b[:0], value, options.flags())
	out := append([]byte(nil), buf...)

	pb.b = buf[:0]
	bytesPool.Put(pb)

	return out, err
}
