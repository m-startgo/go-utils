package m_cron

import (
	"errors"
	"fmt"

	"github.com/robfig/cron/v3"
)

/*
Package m_cron 提供简单的定时任务封装。

用法示例：

c, err := m_cron.New(m_cron.CronOption{
		Func: func() { fmt.Println("Hello") },
		Spec: "0 18 5 3,9,15,21,27 * ? ",
})
if err != nil {
		log.Fatal(err)
}
defer c.Stop() // 程序退出前停止定时器
*/

// CronOption 用于配置定时任务。
type CronOption struct {
	Func func() // 定时执行的函数
	Spec string // cron 表达式
}

// New 创建并启动一个定时任务，返回 cron 实例和错误。
// opt: 定时任务配置选项
func New(opt CronOption) (*cron.Cron, error) {
	if opt.Func == nil {
		return nil, errors.New("参数 Func 不能为空")
	}
	if opt.Spec == "" {
		return nil, errors.New("参数 Spec 不能为空")
	}
	crontab := cron.New(cron.WithSeconds())
	_, err := crontab.AddFunc(opt.Spec, opt.Func)
	if err != nil {
		return nil, fmt.Errorf("定时任务创建失败，Spec 不合法: %w", err)
	}
	crontab.Start()
	return crontab, nil
}
