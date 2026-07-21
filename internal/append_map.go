package internal

// AppendStringBoolMap appends the JSON representation of values to dst.
func AppendStringBoolMap(dst []byte, values map[string]bool) []byte {
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
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringIntMap appends the JSON representation of values to dst.
func AppendStringIntMap[T signedInteger](dst []byte, values map[string]T) []byte {
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
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringUintMap appends the JSON representation of values to dst.
func AppendStringUintMap[T unsignedInteger](dst []byte, values map[string]T) []byte {
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
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringFloat32Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat32Map(dst []byte, values map[string]float32) ([]byte, error) {
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
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringFloat64Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat64Map(dst []byte, values map[string]float64) ([]byte, error) {
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
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringStringMap appends the JSON representation of values to dst.
func AppendStringStringMap(dst []byte, values map[string]string) []byte {
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
		dst = AppendString(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntBoolMap appends the JSON representation of values to dst.
func AppendIntBoolMap[K signedInteger](dst []byte, values map[K]bool) []byte {
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
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntIntMap appends the JSON representation of values to dst.
func AppendIntIntMap[K signedInteger, V signedInteger](dst []byte, values map[K]V) []byte {
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
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntUintMap appends the JSON representation of values to dst.
func AppendIntUintMap[K signedInteger, V unsignedInteger](dst []byte, values map[K]V) []byte {
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
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntFloat32Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendIntFloat32Map[K signedInteger](dst []byte, values map[K]float32) ([]byte, error) {
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
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendIntFloat64Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendIntFloat64Map[K signedInteger](dst []byte, values map[K]float64) ([]byte, error) {
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
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendIntStringMap appends the JSON representation of values to dst.
func AppendIntStringMap[K signedInteger](dst []byte, values map[K]string) []byte {
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
		dst = AppendString(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendUintBoolMap appends the JSON representation of values to dst.
func AppendUintBoolMap[K unsignedInteger](dst []byte, values map[K]bool) []byte {
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
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendUintIntMap appends the JSON representation of values to dst.
func AppendUintIntMap[K unsignedInteger, V signedInteger](dst []byte, values map[K]V) []byte {
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
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendUintUintMap appends the JSON representation of values to dst.
func AppendUintUintMap[K unsignedInteger, V unsignedInteger](dst []byte, values map[K]V) []byte {
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
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendUintFloat32Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendUintFloat32Map[K unsignedInteger](dst []byte, values map[K]float32) ([]byte, error) {
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
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendUintFloat64Map appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendUintFloat64Map[K unsignedInteger](dst []byte, values map[K]float64) ([]byte, error) {
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
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendUintStringMap appends the JSON representation of values to dst.
func AppendUintStringMap[K unsignedInteger](dst []byte, values map[K]string) []byte {
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
		dst = AppendString(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringBoolMapHTML appends the JSON representation of values to dst.
func AppendStringBoolMapHTML(dst []byte, values map[string]bool) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')
		dst = AppendBool(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringIntMapHTML appends the JSON representation of values to dst.
func AppendStringIntMapHTML[T signedInteger](dst []byte, values map[string]T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')
		dst = AppendInt(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringUintMapHTML appends the JSON representation of values to dst.
func AppendStringUintMapHTML[T unsignedInteger](dst []byte, values map[string]T) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')
		dst = AppendUint(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendStringFloat32MapHTML appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat32MapHTML(dst []byte, values map[string]float32) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')

		var err error
		dst, err = AppendFloat32(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringFloat64MapHTML appends the JSON representation of values to dst.
// It returns an error if any value is NaN or infinite.
func AppendStringFloat64MapHTML(dst []byte, values map[string]float64) ([]byte, error) {
	if values == nil {
		return AppendNull(dst), nil
	}
	if len(values) == 0 {
		return append(dst, "{}"...), nil
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')

		var err error
		dst, err = AppendFloat64(dst, value)
		if err != nil {
			return dst, err
		}
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst, nil
}

// AppendStringStringMapHTML appends the JSON representation of values to dst.
func AppendStringStringMapHTML(dst []byte, values map[string]string) []byte {
	if values == nil {
		return AppendNull(dst)
	}
	if len(values) == 0 {
		return append(dst, "{}"...)
	}

	dst = append(dst, '{')
	for key, value := range values {
		dst = AppendStringHTML(dst, key)
		dst = append(dst, ':')
		dst = AppendStringHTML(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendIntBoolMapHTML appends the JSON representation of values to dst.
func AppendIntStringMapHTML[K signedInteger](dst []byte, values map[K]string) []byte {
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
		dst = AppendStringHTML(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}

// AppendUintBoolMapHTML appends the JSON representation of values to dst.
func AppendUintStringMapHTML[K unsignedInteger](dst []byte, values map[K]string) []byte {
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
		dst = AppendStringHTML(dst, value)
		dst = append(dst, ',')
	}
	dst[len(dst)-1] = '}'

	return dst
}
