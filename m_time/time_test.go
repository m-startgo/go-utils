package m_time

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t1 := New()
	if t1 == nil {
		t.Error("New() should not return nil")
	}
}

func TestNewFromTime(t *testing.T) {
	t1 := time.Now()
	t2 := NewFromTime(t1)
	if t2 == nil {
		t.Error("NewFromTime() should not return nil")
	}
	if t2.Time() != t1 {
		t.Error("NewFromTime() should return the same time")
	}
}

func TestNewFromString(t *testing.T) {
	t1, err := NewFromString("2020-01-01")
	if err != nil {
		t.Error("NewFromString() should not return an error")
	}
	if t1 == nil {
		t.Error("NewFromString() should not return nil")
	}

	_, err = NewFromString("invalid")
	if err == nil {
		t.Error("NewFromString() should return an error for invalid string")
	}
}

func TestFormat(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:00:00")
	formatted := t1.Format("2006-01-02 15:04:05")
	if formatted != "2020-01-01 12:00:00" {
		t.Error("Format() should return the correct format")
	}
}

func TestAdd(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:00:00")
	t2 := t1.Add(time.Hour)
	if t2.Format("2006-01-02 15:04:05") != "2020-01-01 13:00:00" {
		t.Error("Add() should add the correct duration")
	}
}

func TestSubtract(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:00:00")
	t2 := t1.Subtract(time.Hour)
	if t2.Format("2006-01-02 15:04:05") != "2020-01-01 11:00:00" {
		t.Error("Subtract() should subtract the correct duration")
	}
}

func TestStartOf(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:30:30")
	t2 := t1.StartOf("day")
	if t2.Format("2006-01-02 15:04:05") != "2020-01-01 00:00:00" {
		t.Error("StartOf() should return the start of the day")
	}
}

func TestEndOf(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:30:30")
	t2 := t1.EndOf("day")
	if t2.Format("2006-01-02 15:04:05") != "2020-01-01 23:59:59" {
		t.Error("EndOf() should return the end of the day")
	}
}

func TestDiff(t *testing.T) {
	t1, _ := NewFromString("2020-01-02 12:00:00")
	t2, _ := NewFromString("2020-01-01 12:00:00")
	diff := t1.Diff(t2, "day")
	if diff.String() != "1" {
		t.Error("Diff() should return the correct difference")
	}
}

func TestUnix(t *testing.T) {
	t1, _ := NewFromString("2020-01-01 12:00:00")
	unix := t1.Unix()
	if unix != 1577880000 {
		t.Error("Unix() should return the correct Unix timestamp")
	}
}

func TestDaysInMonth(t *testing.T) {
	t1, _ := NewFromString("2020-01-01")
	days := t1.DaysInMonth()
	if days != 31 {
		t.Error("DaysInMonth() should return the correct number of days")
	}
}