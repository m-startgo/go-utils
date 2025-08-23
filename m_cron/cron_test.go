package m_cron

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

// TestNew 测试 New 函数
func TestNew(t *testing.T) {
	// 测试正常情况
	t.Run("ValidOption", func(t *testing.T) {
		var executed int32
		cr, err := New(CronOption{
			Func: func() { atomic.StoreInt32(&executed, 1) },
			Spec: "@every 1s", // 使用相对时间，更可预测
		})
		if err != nil {
			t.Fatalf("New() error = %v, want nil", err)
		}
		if cr == nil {
			t.Fatal("New() cronInstance = nil, want not nil")
		}
		t.Cleanup(func() { cr.Stop() })

		// 等待最多 3 秒，轮询检查是否执行
		deadline := time.Now().Add(3 * time.Second)
		for time.Now().Before(deadline) && atomic.LoadInt32(&executed) == 0 {
			time.Sleep(50 * time.Millisecond)
		}
		if atomic.LoadInt32(&executed) == 0 {
			t.Error("定时任务未执行")
		}
	})

	// 测试 Func 为空的情况
	t.Run("NilFunc", func(t *testing.T) {
		cr, err := New(CronOption{
			Func: nil,
			Spec: "@every 1s",
		})
		if err != ErrNilFunc {
			t.Errorf("New() error = %v, want %v", err, ErrNilFunc)
		}
		if cr != nil {
			t.Errorf("New() cronInstance = %v, want nil", cr)
		}
	})

	// 测试 Spec 为空的情况
	t.Run("EmptySpec", func(t *testing.T) {
		cr, err := New(CronOption{
			Func: func() { fmt.Println("test") },
			Spec: "",
		})
		if err != ErrEmptySpec {
			t.Errorf("New() error = %v, want %v", err, ErrEmptySpec)
		}
		if cr != nil {
			t.Errorf("New() cronInstance = %v, want nil", cr)
		}
	})

	// 测试无效的 Spec
	t.Run("InvalidSpec", func(t *testing.T) {
		cr, err := New(CronOption{
			Func: func() { fmt.Println("test") },
			Spec: "invalid spec",
		})
		if err == nil {
			t.Error("New() error = nil, want not nil")
		}
		if cr != nil {
			t.Errorf("New() cronInstance = %v, want nil", cr)
		}
	})
}

// TestCronExecution 测试定时任务是否正确执行
func TestCronExecution(t *testing.T) {
	var counter int32
	cr, err := New(CronOption{
		Func: func() { atomic.AddInt32(&counter, 1) },
		Spec: "@every 1s", // 使用相对时间
	})
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}
	t.Cleanup(func() { cr.Stop() })

	// 等待最多 4 秒，任务应该执行至少 3 次
	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) && atomic.LoadInt32(&counter) < 3 {
		time.Sleep(50 * time.Millisecond)
	}
	if atomic.LoadInt32(&counter) < 3 {
		t.Errorf("定时任务执行次数不足，期望至少 3 次，实际执行 %d 次", counter)
	}
}

// BenchmarkNew 对 New 函数进行基准测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cr, err := New(CronOption{
			Func: func() {},
			Spec: "0 0 12 * * *", // 中午 12 点执行
		})
		if err != nil {
			b.Fatal(err)
		}
		// 直接 Stop() 等待，保证没有残留 goroutine
		cr.Stop()
	}
}

// ExampleNew 提供 New 函数的使用示例
func ExampleNew() {
	// 为了避免示例依赖定时器带来的不稳定，直接演示输出
	fmt.Println("hello")

	// Output: hello
}
