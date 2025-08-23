package m_cron

import (
	"fmt"
	"testing"
	"time"
)

// TestNew 测试 New 函数
func TestNew(t *testing.T) {
	// 测试正常情况
	t.Run("ValidOption", func(t *testing.T) {
		var executed bool
		cronInstance, err := New(CronOption{
			Func: func() { executed = true },
			Spec: "*/1 * * * * *", // 每秒执行一次
		})
		if err != nil {
			t.Fatalf("New() error = %v, want nil", err)
		}
		if cronInstance == nil {
			t.Fatal("New() cronInstance = nil, want not nil")
		}

		// 等待一段时间确保任务执行
		time.Sleep(1500 * time.Millisecond)
		if !executed {
			t.Error("定时任务未执行")
		}

		// 停止 cron
		cronInstance.Stop()
	})

	// 测试 Func 为空的情况
	t.Run("NilFunc", func(t *testing.T) {
		cronInstance, err := New(CronOption{
			Func: nil,
			Spec: "*/1 * * * * *",
		})
		if err != ErrNilFunc {
			t.Errorf("New() error = %v, want %v", err, ErrNilFunc)
		}
		if cronInstance != nil {
			t.Errorf("New() cronInstance = %v, want nil", cronInstance)
		}
	})

	// 测试 Spec 为空的情况
	t.Run("EmptySpec", func(t *testing.T) {
		cronInstance, err := New(CronOption{
			Func: func() { fmt.Println("test") },
			Spec: "",
		})
		if err != ErrEmptySpec {
			t.Errorf("New() error = %v, want %v", err, ErrEmptySpec)
		}
		if cronInstance != nil {
			t.Errorf("New() cronInstance = %v, want nil", cronInstance)
		}
	})

	// 测试无效的 Spec
	t.Run("InvalidSpec", func(t *testing.T) {
		cronInstance, err := New(CronOption{
			Func: func() { fmt.Println("test") },
			Spec: "invalid spec",
		})
		if err == nil {
			t.Error("New() error = nil, want not nil")
		}
		if cronInstance != nil {
			t.Errorf("New() cronInstance = %v, want nil", cronInstance)
		}
	})
}

// TestCronExecution 测试定时任务是否正确执行
func TestCronExecution(t *testing.T) {
	counter := 0
	cronInstance, err := New(CronOption{
		Func: func() { counter++ },
		Spec: "*/1 * * * * *", // 每秒执行一次
	})
	if err != nil {
		t.Fatalf("New() error = %v, want nil", err)
	}

	// 等待 3.5 秒，任务应该执行 3 次
	time.Sleep(3500 * time.Millisecond)
	cronInstance.Stop()

	// 检查执行次数
	if counter < 3 {
		t.Errorf("定时任务执行次数不足，期望至少 3 次，实际执行 %d 次", counter)
	}
}

// BenchmarkNew 对 New 函数进行基准测试
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cronInstance, err := New(CronOption{
			Func: func() {},
			Spec: "0 0 12 * * *", // 中午 12 点执行
		})
		if err != nil {
			b.Fatal(err)
		}
		cronInstance.Stop()
	}
}

// ExampleNew 提供 New 函数的使用示例
func ExampleNew() {
	c, err := New(CronOption{
		Func: func() { fmt.Println("hello") },
		Spec: "*/1 * * * * *", // 每秒执行一次，便于示例即时看到输出
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	// 等待任务触发一次
	time.Sleep(1100 * time.Millisecond)
	c.Stop()

	// Output: hello
}
