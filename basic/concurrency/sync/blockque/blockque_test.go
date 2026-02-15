package blockque_test

import (
	"context"
	"study/basic/concurrency/sync/blockque"
	"study/basic/testing/assertion"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试队列阻塞入队
func TestBlockQueue_Offer(t *testing.T) {
	que := blockque.New[int](10)

	// 前 10 个元素入队不会导致队列阻塞
	for i := range 10 {
		r := que.Offer(context.Background(), i)
		assert.True(t, r)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, que.List())

	// 启动一个 goroutine, 在 100ms 后出队一个元素
	go func() {
		// 等待 100ms 后, 从队列中弹出一个元素
		time.Sleep(100 * time.Millisecond)
		val, ok := que.Poll(-1)

		// 确认弹出了队列中第一个元素
		assert.True(t, ok)
		assert.Equal(t, 0, val)
	}()

	start := time.Now()

	// 入队一个新元素, 总体耗时 100ms 以上 (包括等待队列出队)
	r := que.Offer(context.Background(), 10)
	assert.True(t, r)
	assertion.DurationMatch(t, 100*time.Millisecond, time.Since(start))

	// 确认队列的长度和内容
	assert.Equal(t, 10, que.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, que.List())
}

// 测试队列入队失败
func TestBlockQueue_TryOffer(t *testing.T) {
	que := blockque.New[int](10)

	// 前 10 个元素入队不会失败
	for i := range 10 {
		r := que.TryOffer(i)
		assert.True(t, r)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, que.List())

	// 入队第 11 个元素, 由于队列已满, 所以返回 false 表示入队失败
	r := que.TryOffer(10)
	assert.False(t, r)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, que.List())

	// 从队列中删除一个元素
	que.Remove()

	// 再次入队第 11 个元素, 入队成功
	r = que.TryOffer(10)
	assert.True(t, r)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, que.List())
}

// 测试队列入队超时
func TestBlockQueue_OfferWithTimeout(t *testing.T) {
	que := blockque.New[int](10)

	start := time.Now()

	// 前 10 个元素入队不会超时
	for i := range 10 {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

		r := que.Offer(ctx, i)
		cancel()

		assert.True(t, r)
	}
	// 确认前 10 个元素入队无需等待
	assert.Less(t, time.Since(start), 10*time.Millisecond)
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, que.List())

	start = time.Now()

	// 入队第 11 个元素, 由于队列已满, 等待 100ms 后超时, 入队失败
	ctx, cancel := context.WithTimeout(context.Background(), 109*time.Millisecond)
	r := que.Offer(ctx, 10)
	cancel()

	// 确认第 11 个元素入队失败
	assert.False(t, r)
	// 确认入队第 11 元素时等待了 100ms 后超时失败
	assertion.DurationMatch(t, 100*time.Millisecond, time.Since(start))
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, que.List())

	// 从队列中删除一个元素
	que.Remove()

	start = time.Now()

	// 再次入队第 11 个元素, 由于已经从队列取出一个元素, 则本次入队成功
	r = que.Offer(context.Background(), 10)

	// 确认第 11 个元素再次入队成功
	assert.True(t, r)
	// 确认入队第 11 元素时未发生等待
	assert.Less(t, time.Since(start), 10*time.Millisecond)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, que.List())
}
