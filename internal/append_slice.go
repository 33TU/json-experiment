package internal

// AppendBoolSlice appends the JSON representation of values to dst.
func AppendBoolSlice(dst []byte, values []bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendIntSlice appends the JSON representation of values to dst.
func AppendIntSlice[T signedInteger](dst []byte, values []T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendUintSlice appends the JSON representation of values to dst.
func AppendUintSlice[T unsignedInteger](dst []byte, values []T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendFloat32Slice appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendFloat32Slice(dst []byte, values []float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "[]"...), nil
	}

	dst = append(dst, '[')
	for _, value := range values {
		var err error
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst, nil
}

// AppendFloat64Slice appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendFloat64Slice(dst []byte, values []float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "[]"...), nil
	}

	dst = append(dst, '[')
	for _, value := range values {
		var err error
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst, nil
}

// AppendStringSlice appends the JSON representation of values to dst.
func AppendStringSlice(dst []byte, values []string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendString(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}

// AppendStringSliceHTML appends the HTML-safe JSON representation of values to dst.
func AppendStringSliceHTML(dst []byte, values []string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "[]"...)
	}

	dst = append(dst, '[')
	for _, value := range values {
		dst = AppendStringHTML(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = ']'

	return dst
}
