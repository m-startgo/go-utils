package m_cycle

import (
	"sync"
	"sync/atomic"
	"time"
)

/*
cy := m_cycle.New(m_cycle.Opt{Func: myFunc, SleepTime: time.Second})
cy.Start()
time.Sleep(5 * time.Second)
cy.EditTime(2 * time.Second) // 动态修改为2秒

*/

type Cycle struct {
	Task     func()
	Interval time.Duration

	status int32
	stopCh chan struct{}
	editCh chan time.Duration
	once   sync.Once
}

type Options struct {
	Task     func()
	Interval time.Duration
}

func New(opt Options) *Cycle {
	c := &Cycle{
		Task:     opt.Task,
		Interval: opt.Interval,
		status:   1,
		stopCh:   make(chan struct{}),
		editCh:   make(chan time.Duration, 1), // 缓冲防阻塞
	}
	if c.Task == nil {
		c.Task = func() {}
	}
	return c
}

func (c *Cycle) End() *Cycle {
	c.once.Do(func() {
		if atomic.CompareAndSwapInt32(&c.status, 1, 2) {
			close(c.stopCh)
		}
	})
	return c
}

func (c *Cycle) Start() *Cycle {
	c.Task()
	go func() {
		ticker := time.NewTicker(c.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-c.stopCh:
				return
			case d := <-c.editCh:
				ticker.Stop()
				ticker = time.NewTicker(d)
			case <-ticker.C:
				c.Task()
			}
		}
	}()
	return c
}

// SetInterval 动态修改定时器间隔
func (c *Cycle) SetInterval(d time.Duration) {
	c.Interval = d
	select {
	case c.editCh <- d:
	default:
	}
}
