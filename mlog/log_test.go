package mlog

import "testing"

var Log *Logger

// go test -v -run TestNewLog
func TestNewLog(t *testing.T) {
	Log = New(Config{
		Path: "./mo7-logs",
		Name: "m",
	})
	Log.Info("this is info")
	Log.Warn("this is warn")
	Log.Error("this is error")
	Log.Debug("this is debug")
}

// go test -v -run TestClearLog
func TestClearLog(t *testing.T) {
	Log = New(Config{
		Path: "./mo7-logs",
		Name: "m-logtest",
	})

	Log.Clear(ClearOpt{
		Type:   []string{"debug", "warn"},
		Before: 7,
	})
}
