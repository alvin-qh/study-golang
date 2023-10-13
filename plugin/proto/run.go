package proto

type Runnable interface {
	Run(args string) int
}
