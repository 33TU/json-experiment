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
	return marshalInterface(dst, value)
}

func Marshal(value any) ([]byte, error) {
	pb := bytesPool.Get().(*pooledBuffer)

	buf, err := marshalInterface(pb.b[:0], value)
	out := append([]byte(nil), buf...)

	pb.b = buf[:0]
	bytesPool.Put(pb)

	return out, err
}
