package internal

// AppendStringBoolSliceMap appends the JSON representation of values to dst.
func AppendStringBoolSliceMap(dst []byte, values map[string][]bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')
		dst = AppendBoolSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringIntSliceMap appends the JSON representation of values to dst.
func AppendStringIntSliceMap[T signedInteger](dst []byte, values map[string][]T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')
		dst = AppendIntSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringUintSliceMap appends the JSON representation of values to dst.
func AppendStringUintSliceMap[T unsignedInteger](dst []byte, values map[string][]T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')
		dst = AppendUintSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringFloat32SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat32SliceMap(dst []byte, values map[string][]float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')

		var err error
		dst, err = AppendFloat32Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringFloat64SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat64SliceMap(dst []byte, values map[string][]float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')

		var err error
		dst, err = AppendFloat64Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringStringSliceMap appends the JSON representation of values to dst.
func AppendStringStringSliceMap(dst []byte, values map[string][]string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendString(dst, key)
		dst = append(dst, ':')
		dst = AppendStringSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntBoolSliceMap appends the JSON representation of values to dst.
func AppendIntBoolSliceMap[K signedInteger](dst []byte, values map[K][]bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendBoolSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendIntIntSliceMap appends the JSON representation of values to dst.
func AppendIntIntSliceMap[K signedInteger, V signedInteger](dst []byte, values map[K][]V) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendIntSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendIntUintSliceMap appends the JSON representation of values to dst.
func AppendIntUintSliceMap[K signedInteger, V unsignedInteger](dst []byte, values map[K][]V) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendUintSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendIntFloat32SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendIntFloat32SliceMap[K signedInteger](dst []byte, values map[K][]float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		var err error
		dst, err = AppendFloat32Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst, nil
}

// AppendIntFloat64SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendIntFloat64SliceMap[K signedInteger](dst []byte, values map[K][]float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		var err error
		dst, err = AppendFloat64Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst, nil
}

// AppendIntStringSliceMap appends the JSON representation of values to dst.
func AppendIntStringSliceMap[K signedInteger](dst []byte, values map[K][]string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendInt(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendStringSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendUintBoolSliceMap appends the JSON representation of values to dst.
func AppendUintBoolSliceMap[K unsignedInteger](dst []byte, values map[K][]bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendBoolSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendUintIntSliceMap appends the JSON representation of values to dst.
func AppendUintIntSliceMap[K unsignedInteger, V signedInteger](dst []byte, values map[K][]V) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendIntSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendUintUintSliceMap appends the JSON representation of values to dst.
func AppendUintUintSliceMap[K unsignedInteger, V unsignedInteger](dst []byte, values map[K][]V) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendUintSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}

// AppendUintFloat32SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendUintFloat32SliceMap[K unsignedInteger](dst []byte, values map[K][]float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		var err error
		dst, err = AppendFloat32Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst, nil
}

// AppendUintFloat64SliceMap appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendUintFloat64SliceMap[K unsignedInteger](dst []byte, values map[K][]float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		var err error
		dst, err = AppendFloat64Slice(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst, nil
}

// AppendUintStringSliceMap appends the JSON representation of values to dst.
func AppendUintStringSliceMap[K unsignedInteger](dst []byte, values map[K][]string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}
	dst = append(dst, '{')
	for key, value := range values {
		dst = append(dst, '"')
		dst = AppendUint(dst, key)
		dst = append(dst, '"', ':')
		dst = AppendStringSlice(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'
	return dst
}
