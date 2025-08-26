package m_cycle

import (
	"sync/atomic"
	"testing"
	"time"
)

// go test -v -run TestCycle_Basic
func TestCycle_Basic(t *testing.T) {
	var count int32
	cy := New(Options{
		Task: func() {
			atomic.AddInt32(&count, 1)
		},
		Interval: 50 * time.Millisecond,
	})
	cy.Start()
	time.Sleep(120 * time.Millisecond)
	cy.End()
	val := atomic.LoadInt32(&count)
	if val < 2 {
		t.Errorf("期望至少执行2次，实际为 %d", val)
	}
}

// go test -v -run TestCycle_SetInterval

func TestCycle_SetInterval(t *testing.T) {
	var count int32
	cy := New(Options{
		Task: func() {
			atomic.AddInt32(&count, 1)
		},
		Interval: 100 * time.Millisecond,
	})
	cy.Start()
	time.Sleep(120 * time.Millisecond)
	cy.SetInterval(20 * time.Millisecond)
	time.Sleep(70 * time.Millisecond)
	cy.End()
	val := atomic.LoadInt32(&count)
	if val < 4 {
		t.Errorf("修改间隔后期望至少执行4次，实际为 %d", val)
	}
}

// go test -v -run TestCycle_EndIdempotent

func TestCycle_EndIdempotent(t *testing.T) {
	cy := New(Options{
		Task:     func() {},
		Interval: 10 * time.Millisecond,
	})
	cy.Start()
	cy.End()
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("End 应该是幂等的，不应发生 panic")
		}
	}()
	cy.End()
}

// go test -v -run TestCycle_NilTask

func TestCycle_NilTask(t *testing.T) {
	cy := New(Options{
		Task:     nil,
		Interval: 10 * time.Millisecond,
	})
	cy.Start()
	time.Sleep(20 * time.Millisecond)
	cy.End()
}

// go test -v -run TestCycle_ImmediateEnd

func TestCycle_ImmediateEnd(t *testing.T) {
	var count int32
	cy := New(Options{
		Task: func() {
			atomic.AddInt32(&count, 1)
		},
		Interval: 10 * time.Millisecond,
	})
	cy.End()
	cy.Start()
	time.Sleep(20 * time.Millisecond)
	val := atomic.LoadInt32(&count)
	if val > 1 {
		t.Errorf("期望最多执行1次，实际为 %d", val)
	}
}
