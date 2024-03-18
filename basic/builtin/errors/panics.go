package errors

type List struct {
	slice []string
}

func (l *List) Append(s string) {
	l.slice = append(l.slice, s)
}
