package pool

import (
	"os"
	"study/basic/logs"
	"sync"
	"sync/atomic"
)

var (
	lastTaskId atomic.Int64
	logger     *logs.Logger
)

// 初始化任务池
func init() {
	// 初始化日志实例
	logger = logs.New()
	logger.AddNewAppender(os.Stdout, logs.LEVEL_DEBUG, logs.Ldate|logs.Ltime|logs.Lshortfile)
}

// 任务处理函数类型
type TaskHandler[T, R any] func(arg T) (R, error)

// 任务类型
type Task[T, R any] struct {
	Id        int64             // 任务 id, 用于任务追踪
	handler   TaskHandler[T, R] // 执行任务的函数
	Argument  T                 // 任务参数
	onSuccess func(R)           // 任务成功的回调函数
	onError   func(error)       // 任务失败的回调函数
}

// 任务池类型
type TaskPool[T, R any] struct {
	ch  chan *Task[T, R] // 输送任务实例的通道
	wg  sync.WaitGroup   // 等待任务结束的等待组
	mux sync.Mutex
}

// 创建任务池实例
func NewTaskPool[T, R any](size int) *TaskPool[T, R] {
	// 创建任务池实例
	pool := TaskPool[T, R]{
		ch: make(chan *Task[T, R], size),
	}

	// 启动指定数量的 goroutine, 执行异步任务
	for i := 0; i < size; i++ {
		pool.wg.Add(1)

		go func() {
			defer pool.wg.Done()

			logger.Debug("worker %d starting", i)

			// 在循环中获取任务并执行
			for {
				ch := pool.channel()
				if ch == nil {
					logger.Debug("worker %d canceled", i)
					return
				}

				// 获取一个任务实例
				t, ok := <-ch
				if !ok {
					logger.Debug("worker %d canceled", i)
					return
				}
				logger.Debug("new task incoming, id: %d", t.Id)

				// 执行任务, 并返回结果
				r, err := t.handler(t.Argument)

				// 根据任务执行是否成功调用不同的回调函数
				if err == nil {
					t.onSuccess(r)
				} else {
					t.onError(err)
				}
			}
		}()
	}

	return &pool
}

func (p *TaskPool[T, R]) channel() chan *Task[T, R] {
	p.mux.Lock()
	defer p.mux.Unlock()
	ch := p.ch
	return ch
}

// 执行一个任务
func (p *TaskPool[T, R]) Worker(handler TaskHandler[T, R]) func(T, func(R), func(error)) {
	return func(arg T, onSuccess func(R), onError func(error)) {
		// 创建任务实例
		t := &Task[T, R]{
			Id:        lastTaskId.Add(1),
			handler:   handler,
			Argument:  arg,
			onSuccess: onSuccess,
			onError:   onError,
		}
		logger.Debug("new task created, id: %v", t.Id)

		// 将任务实例发送到通道
		p.ch <- t
	}
}

func (p *TaskPool[T, R]) close() bool {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.ch == nil {
		return false
	}

	close(p.ch)
	p.ch = nil

	return true
}

// 关闭任务池
func (p *TaskPool[T, R]) Close() {
	p.close()
}

// 关闭任务池并等待所有任务结束
func (p *TaskPool[T, R]) CloseAndWait() {
	if p.close() {
		p.wg.Wait()
	}
}
