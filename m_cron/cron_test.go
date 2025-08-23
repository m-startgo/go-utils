package m_cron

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// TestNew_Errors 验证参数校验与解析错误
func TestNew_Errors(t *testing.T) {
	// Func 为空
	if _, err := New(Options{}); err == nil {
		t.Fatal("期望 Func 为空时报错")
	}

	// Spec 为空
	if _, err := New(Options{Func: func() {}}); err == nil {
		t.Fatal("期望 Spec 为空时报错")
	}

	// 非法表达式
	if _, err := New(Options{Func: func() {}, Spec: "!!! invalid !!!"}); err == nil {
		t.Fatal("期望非法表达式时报错")
	}
}

// TestJob_RunStopWait 验证任务能执行，Stop 会等待正在运行的任务完成
func TestJob_RunStopWait(t *testing.T) {
	var ran int32
	started := make(chan struct{}, 1)

	fn := func() {
		// 标记开始
		started <- struct{}{}
		// 模拟耗时任务
		time.Sleep(150 * time.Millisecond)
		atomic.AddInt32(&ran, 1)
	}

	// 使用 descriptor 表达式以便快速触发：每 200ms 执行一次
	job, err := New(Options{
		Func:        fn,
		Spec:        "@every 200ms",
		WithSeconds: false,
	})
	if err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}
	// 确保至少触发一次
	select {
	case <-started:
	case <-time.After(1 * time.Second):
		job.Stop()
		t.Fatal("任务未在预期时间内开始执行")
	}

	// 在任务正在运行时调用 Stop，Stop 应等待任务完成
	job.Stop()

	// 停止后应至少完成一次执行
	if cnt := atomic.LoadInt32(&ran); cnt < 1 {
		t.Fatalf("期望至少运行 1 次，实际 %d 次", cnt)
	}
}

// TestJob_RemoveStopsScheduling 验证 Remove 能移除后续调度（短时间内）
func TestJob_RemoveStopsScheduling(t *testing.T) {
	var cnt int32
	fn := func() {
		atomic.AddInt32(&cnt, 1)
	}

	job, err := New(Options{
		Func:        fn,
		Spec:        "@every 150ms",
		WithSeconds: false,
	})
	if err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}
	// 等待第一次触发
	waitUntil := time.After(1 * time.Second)
	for {
		if atomic.LoadInt32(&cnt) >= 1 {
			break
		}
		select {
		case <-waitUntil:
			job.Stop()
			t.Fatal("未在预期时间内收到第一次触发")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	// 记录当前次数并移除
	before := atomic.LoadInt32(&cnt)
	job.Remove()

	// 等待一段时间确认没有新触发（间隔为 150ms，等待 3 次间隔）
	time.Sleep(500 * time.Millisecond)
	after := atomic.LoadInt32(&cnt)

	// 允许极小概率的并发残留执行：只要没有显著增加即可判定成功
	if after > before+1 {
		job.Stop()
		t.Fatalf("Remove 后出现过多次执行，移除失败: before=%d after=%d", before, after)
	}

	// 清理：停止 cron 实例
	job.Stop()
}

// TestWait_WithContextTimeout 验证 Wait 在 ctx 超时或取消时会返回
func TestWait_WithContextTimeout(t *testing.T) {
	var ran int32
	started := make(chan struct{}, 1)

	fn := func() {
		started <- struct{}{}
		// 模拟较长任务
		time.Sleep(300 * time.Millisecond)
		atomic.AddInt32(&ran, 1)
	}

	job, err := New(Options{
		Func:        fn,
		Spec:        "@every 200ms",
		WithSeconds: false,
	})
	if err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}

	// 等待任务开始
	select {
	case <-started:
	case <-time.After(1 * time.Second):
		job.Stop()
		t.Fatal("任务未在预期时间内开始执行")
	}

	// 使用短超时 ctx 调用 Wait，应在超时后返回
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	job.Wait(ctx) // 这里不关心内部是否已完全停止，只需确保 Wait 能在 ctx 超时后返回

	// 清理
	job.Stop()
	if atomic.LoadInt32(&ran) < 1 {
		t.Fatal("期望任务至少执行一次")
	}
}
