package m_cron

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
)

/*
Cron在线表达式生成器: https://cron.ciding.cc/

Package m_cron 提供一个简易的定时任务封装，基于 github.com/robfig/cron/v3
默认使用带秒字段的解析器（cron.WithSeconds()）。
注意：robfig/cron 不完全兼容 Quartz 的一些扩展语法（如 '?'、L、W、# 等）。
若需兼容 Quartz 风格表达式，建议在调用方做预处理或使用专门支持 Quartz 的库。

示例：

	c, err := m_cron.New(CronOption{
			Func: func() { fmt.Println("hello") },
			Spec: "0 0 12 * * *", // 带秒的 robfig/cron 表达式
	})

	if err != nil {
			log.Fatal(err)
	}
	defer c.Stop()
*/
var (
	// ErrNilFunc 表示未提供要执行的函数。
	ErrNilFunc = errors.New("m_cron: Func is nil")
	// ErrEmptySpec 表示未提供 cron 表达式。
	ErrEmptySpec = errors.New("m_cron: Spec is empty")
)

// CronOption 是 New 的配置项。
type CronOption struct {
	Func func() // 定时执行的函数，不能为空
	Spec string // cron 表达式（带秒时为 6 字段），不能为空
}

// New 创建并启动定时任务并返回已启动的 *cron.Cron。
// 当参数不合法或表达式无法解析时返回错误。返回的 cron 需要调用方在适当时机 Stop()。
func New(opt CronOption) (*cron.Cron, error) {
	if opt.Func == nil {
		return nil, ErrNilFunc
	}
	if opt.Spec == "" {
		return nil, ErrEmptySpec
	}

	c := cron.New(cron.WithSeconds())

	if _, err := c.AddFunc(opt.Spec, opt.Func); err != nil {
		return nil, fmt.Errorf("m_cron: invalid spec %q: %w", opt.Spec, err)
	}

	c.Start()
	return c, nil
}
