package mlog

/*

myLog := mlog.New(mlog.Config{
	Path: "./logs",
	Name: "log",
})
myLog.Info("this is info")
myLog.Warn("this is warn")
myLog.Error("this is error")
myLog.Debug("this is debug")

*/

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/m-startgo/go-utils/mfile"
	"github.com/m-startgo/go-utils/mpath"
	"github.com/m-startgo/go-utils/mtime"
)

// Logger 是简单日志器的实例。它负责根据名称、类型和日期将日志写入文件。
type Logger struct {
	Path string // 存放日志的目录，默认为 ./logs
	Name string // 日志名称，默认为 log

	mu sync.Mutex // 写入时的互斥
}

// Config 用于构造 Logger 的配置。
type Config struct {
	Path string
	Name string
}

// ClearOpt 表示清理日志的参数。
type ClearOpt struct {
	Type   []string // 要删除的日志类型，为空则删除全部类型
	Before int      // 保留多少天之前的日志会被删除（天数），为空或<=0 则默认90
}

// validTypes 是允许的日志类型集合。
var validTypes = map[string]struct{}{
	"info":  {},
	"warn":  {},
	"error": {},
	"debug": {},
}

// New 根据传入的配置创建并返回一个 Logger 实例。
// 如果 Path 为空则使用 ./logs；如果 Name 为空则使用 log。
func New(cfg Config) *Logger {
	p := cfg.Path
	if strings.TrimSpace(p) == "" {
		p = "./logs"
	}
	n := cfg.Name
	if strings.TrimSpace(n) == "" {
		n = "log"
	}
	// 若目录不存在，mfile.Write/Append 会尝试创建；但这里提前创建以便行为更明确
	if !mpath.IsExist(p) {
		_ = os.MkdirAll(p, 0o755)
	}
	return &Logger{Path: p, Name: n}
}

// formatLine 将要写入日志文件的单行格式化。
func formatLine(level string, v ...any) string {
	ts := mtime.Now().FormatDefault()
	msg := fmt.Sprintln(v...)
	// 去掉 Sprintln 带来的末尾换行，统一以单行记录（末尾保留换行用于文件分隔）
	msg = strings.TrimRight(msg, "\n")
	return fmt.Sprintf("[%s] [%s] %s\n", ts, strings.ToUpper(level), msg)
}

// logToFile 将格式化后的行追加到对应类型的日志文件中。
func (l *Logger) logToFile(typ string, line string) error {
	// guard against nil receiver
	if l == nil {
		return fmt.Errorf("err:mlog.logToFile|nil logger")
	}

	if typ == "" {
		typ = "info"
	}

	// protect against empty Path/Name
	path := strings.TrimSpace(l.Path)
	if path == "" {
		path = "./logs"
	}
	name := strings.TrimSpace(l.Name)
	if name == "" {
		name = "log"
	}

	// ensure directory exists
	if !mpath.IsExist(path) {
		_ = os.MkdirAll(path, 0o755)
	}

	date := mtime.FromTime(time.Now()).Format("YYYY-MM-DD")
	filename := fmt.Sprintf("%s-%s-%s.log", name, typ, date)
	fp := filepath.Join(path, filename)
	l.mu.Lock()
	defer l.mu.Unlock()
	return mfile.Append(fp, line)
}

// Info 写一条 info 级别的日志。
func (l *Logger) Info(v ...any) error {
	return l.logToFile("info", formatLine("info", v...))
}

// Warn 写一条 warn 级别的日志。
func (l *Logger) Warn(v ...any) error {
	return l.logToFile("warn", formatLine("warn", v...))
}

// Error 写一条 error 级别的日志。
func (l *Logger) Error(v ...any) error {
	return l.logToFile("error", formatLine("error", v...))
}

// Debug 写一条 debug 级别的日志。
func (l *Logger) Debug(v ...any) error {
	return l.logToFile("debug", formatLine("debug", v...))
}
