package m_cron

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
)

// Cron在线表达式生成器: https://cron.ciding.cc/

// Job 封装 cron 实例和 entry id，提供 Stop/Remove/Wait
type Job struct {
	C      *cron.Cron
	Entry  cron.EntryID
	cancel context.CancelFunc
}

// Options 定时任务参数（简化）
type Options struct {
	Func          func()
	Spec          string
	BlockMain     bool            // 是否阻塞主线程并在信号到来时停止
	Ctx           context.Context // 可选上下文，用于外部取消
	Location      *time.Location  // 默认 nil（即 time.Local）
	WithSeconds   bool            // 默认 false（如果需要秒级表达式请设为 true）
	RecoverLogger cron.Logger     // 默认 nil（即 cron.DefaultLogger）
}

// New 创建并启动定时任务（简化版）
// 使用示例：
//
//	job, err := m_cron.New(m_cron.Options{
//	    Func:        func(){ fmt.Println("tick") },
//	    Spec:        "0 * * * * *", // 含秒表达式需 WithSeconds=true
//	    WithSeconds: true,
//	})
func New(opt Options) (*Job, error) {
	if opt.Func == nil {
		return nil, fmt.Errorf("参数 Func 不能为空")
	}
	opt.Spec = strings.TrimSpace(opt.Spec)
	if opt.Spec == "" {
		return nil, fmt.Errorf("参数 Spec 不能为空")
	}

	// 校验表达式：根据 WithSeconds 简单校验
	var parser cron.Parser
	if opt.WithSeconds {
		parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	} else {
		parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}
	if _, err := parser.Parse(opt.Spec); err != nil {
		return nil, fmt.Errorf("表达式解析失败: %w", err)
	}

	loc := opt.Location
	if loc == nil {
		loc = time.Local
	}
	rl := opt.RecoverLogger
	if rl == nil {
		rl = cron.DefaultLogger
	}

	// 为兼容秒级表达式，始终启用 WithSeconds（安全且简单）
	c := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(loc),
		cron.WithChain(cron.Recover(rl)),
	)

	id, err := c.AddFunc(opt.Spec, opt.Func)
	if err != nil {
		return nil, fmt.Errorf("添加任务失败: %w", err)
	}

	c.Start()

	// 上下文用于外部取消
	var ctx context.Context
	var cancel context.CancelFunc
	if opt.Ctx != nil {
		ctx, cancel = context.WithCancel(opt.Ctx)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	_ = ctx // 保留以便将来使用 / 扩展

	job := &Job{
		C:      c,
		Entry:  id,
		cancel: cancel,
	}

	// 可选阻塞主线程并在信号到来时停止
	if opt.BlockMain {
		sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()
		<-sigCtx.Done()
		job.Stop()
	}

	return job, nil
}

// Stop 优雅停止并等待任务完成
func (j *Job) Stop() {
	if j == nil || j.C == nil {
		return
	}
	if j.cancel != nil {
		j.cancel()
	}
	done := j.C.Stop()
	<-done.Done()
}

// Remove 从调度器移除任务（若要确保任务已停止，请先 Stop）
func (j *Job) Remove() {
	if j == nil || j.C == nil {
		return
	}
	j.C.Remove(j.Entry)
}

// Wait 在 ctx 下等待任务停止（ctx 为 nil 则立即返回）
func (j *Job) Wait(ctx context.Context) {
	if j == nil || j.C == nil || ctx == nil {
		return
	}
	done := j.C.Stop()
	select {
	case <-done.Done():
	case <-ctx.Done():
	}
}
