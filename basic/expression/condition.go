package expression

// 根据条件返回不同值
func If[T any](cond bool, first T, second T) T {
	if cond {
		return first
	}
	return second
}

// 保存 `err` 参数必须为 `nil` 的函数, 返回 `v` 参数
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
