

package context

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// 在使用 Value Context 时, Key 的类型推荐使用自定义类型
type ContextKey string

var (
	ErrInvalidContext = errors.New("invalid context")
)

// 创建一个 Value Context
//
// 可以通过 context.Context 携带所需的值 (以 Key/Value) 形式, 传递到下层的 goroutine 协程中
func CreateContext() (context.Context, chan string) {
	ch := make(chan string)                                                    // 创建一个 chan 对象
	return context.WithValue(context.Background(), ContextKey("chan"), ch), ch // 创建一个 Value Context, 包含 chan 对象
}

// 创建一个 Cancel Context
//
// 可以通过 context.Context 的 Done 函数判断是否已被标志为已取消, 所有下层的 goroutine 协程都会被取消
func CreateCancelContext() (context.Context, context.CancelFunc, chan string) {
	ctx, ch := CreateContext()
	ctx, cancel := context.WithCancel(ctx) // 创建 Cancel Context, 返回 Context 对象和 cancel 函数
	return ctx, cancel, ch
}

// 创建一个 Timeout Context
//
// 可以通过 `context.Context` 的 `Done` 函数判断是否已被标志为已取消或超时, 所有下层的 goroutine 协程都会被取消
// Timeout Context 同时具备 Cancel Context 和超时取消的特性
func CreateTimeoutContext(timeout time.Duration) (context.Context, context.CancelFunc, chan string) {
	ctx, ch := CreateContext()
	ctx, cancel := context.WithTimeout(ctx, timeout) // 创建 Timeout Context, 返回 Context 对象和 cancel 函数
	return ctx, cancel, ch
}

// 测试 Cancel Context
//
// 调用 `context.WithCancel` 可以创建一个带有取消功能的 `Context`, 该函数返回两个结果: Context 和 CancelFunc
//  1. Context 的 `Done` 函数返回一个 `chan`, 如果接收到数据 (为空结构体值 `struct{}{}`), 则表示协程需要被取消
//  2. CancelFunc 如果被调用, 则 1. 中的 `Done` 函数返回的 `chan` 会接收到数据, 表示协程需要被取消
func CancelContextHandler(ctx context.Context) {
	if ctx.Done() == nil {
		panic(ErrInvalidContext)
	}

	// 从 Value Context 中获取 chan 类型值
	ch := ctx.Value(ContextKey("chan")).(chan string)

	// 根据 chan 的状态绝对下一步操作
	select {
	case <-ctx.Done(): // Cancel Context 的 cancel 函数被调用, 需要取消协程
		fmt.Printf("Work canceled because: %v\n", ctx.Err())
		ch <- "Canceled"
	case result := <-ch: // 接收到 chan 传入的值, 然后正常结束当前协程
		fmt.Printf("Work completed, the result is: %v\n", result)
		ch <- "Completed"
	}
}

// 测试 Timeout Context
//
// 调用 `context.WithTimeout` 可以创建一个带有超时取消功能的 `Context`, 该函数返回两个结果, Context 和 ContextFunc
//
//  1. Context 的 `Done` 函数返回一个 `chan`, 如果接收到数据 (为空结构体值 `struct{}{}`), 则表示协程需要被取消
//  2. CancelFunc 如果被调用, 则 1. 中的 `Done` 函数返回的 `chan` 会接收到数据, 表示协程需要被取消
//  3. 当 Timeout Context 的超时时间到达后, 则 2. 中描述的情况也会发生
func TimeoutHandler(ctx context.Context, executeTime time.Duration) {
	// 从 Value Context 中获取 chan 类型值
	ch := ctx.Value(ContextKey("chan")).(chan string)

	select {
	case <-ctx.Done(): // 收到协程结束通知 (cancel 函数被调用) 或超时时间到达, 表示取消或超时结束
		t, ok := ctx.Deadline()
		if ok && time.Since(t) >= 0 {
			fmt.Printf("Work timeout because: %v\n", ctx.Err())
			ch <- "Canceled by timeout"
		} else {
			fmt.Printf("Work timeout because: cancel function called, %v\n", ctx.Err())
			ch <- "Canceled by cancel called"
		}
	case <-time.After(executeTime): // executeTime 参数规定的运行时间到达, 表示正常结束
		fmt.Println("Work completed at time")
		ch <- "Completed by timer"
	case result := <-ch: // 接收到 chan 传入的值后结束当前协程, 表示正常结束
		fmt.Printf("Work completed, the result is: %v\n", result)
		ch <- "Completed by chan"
	}
}
