package pool

type Task struct {
	Id     string
	Data   interface{}
	Result interface{}
	Error  error
}
