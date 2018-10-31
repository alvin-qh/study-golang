package builtin

type RunAfter func(s string) string

func DeferFunc(fn RunAfter) (ret string) {
	defer func() {
		ret = fn(ret)
	}()

	ret = "Hello"
	return
}

func PanicFunc(s string) string {
	if len(s) == 0 {
		panic("Empty")
	}
	return s
}
