package builtin

import "errors"

func WithSafeError(name string) (msg string, err error) {
	if len(name) == 0 {
		err = errors.New("invalid name")
		return
	}
	msg = "Hello " + name
	return
}

func WithPanicError(name string) string {
	msg, err := WithSafeError(name)
	if err != nil {
		panic(err)
	}
	return msg
}
