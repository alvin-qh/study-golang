package conv

// 将 `any` 类型转换为切片类型
//
// 即 `[]any` => `[]T`
func AnyToSlice[T any](v any) ([]T, bool) {
	if s, ok := v.([]T); ok {
		return s, true
	}
	return nil, false
}

// 将指定类型切片转为 `any` 类型切片
//
// 即 `[]T` => `[]any`
func TypedSliceToAnySlice[T any](v []T) []any {
	if len(v) == 0 {
		return nil
	}

	r := make([]any, len(v))
	for i, v := range v {
		r[i] = v
	}
	return r
}

// 将 `any` 类型切片转为指定类型切片
//
// 即 `[]any` => `[]T`
func AnySliceToTypedSlice[T any](v []any) ([]T, bool) {
	if len(v) == 0 {
		return nil, true
	}

	var ok bool
	r := make([]T, len(v))
	for i, v := range v {
		if r[i], ok = v.(T); !ok {
			return r, false
		}
	}
	return r, true
}
