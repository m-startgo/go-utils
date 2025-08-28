package mlog

/*
我想给予标准库封装一个更简洁的日志接口，方便在项目中使用和替换不同的日志实现。

日志文件以 名称-类型-YYYY-MM-DD.log 格式命名，方便按类型和日期查找。
由用户决定日志类型和名称，以及存放目录，如果目录不存在则自动创建。

涉及到时间，路径和文件的处理，使用本地的 mtime 和  mpath 以及 mfile 库

然后 log 可以定义删除

调用方式如下

var myLog = mlog.New({
  Path: "./logs", // 为空 则 默认为 ./logs
	Name: "log", // 为空 则 默认为 log
})

myLog.Info("This is an info message")
myLog.Warn("This is a warning message")
myLog.Error("This is an error message")
myLog.Debug("This is a debug message")


// 可以定义多个删除方法
myLog.Clear({
	Type: []string{
		"warn","debug"
	}, // 为空则 默认删除全部类型
	Before: 7 , // 距离现在的时长 天数，为空 则 默认90天
})
myLog.Clear({
	Type: []string{
		"info","warn"
	}, // 为空则 默认删除全部类型
	Before: 30 , // 距离现在的时长 天数，为空 则 默认90天
})

// 因为日志格式是固定的，所以可以通过名称和类型来删除指定的日志文件
// Type 意思是要清除的日志类型，只能为 "info","warn","error","debug" 这四种

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
	if typ == "" {
		typ = "info"
	}
	date := mtime.FromTime(time.Now()).Format("YYYY-MM-DD")
	filename := fmt.Sprintf("%s-%s-%s.log", l.Name, typ, date)
	fp := filepath.Join(l.Path, filename)
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

// Clear 按配置删除匹配的历史日志文件。
// - Type: 指定要删除的类型（info,warn,error,debug），为空则删除全部类型
// - Before: 距离现在的天数，默认 90 天，删除早于该天数的日志文件
func (l *Logger) Clear(opt ClearOpt) {
	beforeDays := opt.Before
	if beforeDays <= 0 {
		beforeDays = 90
	}
	// 计算截止日期（保留到 beforeDays 天之前的之后文件）
	// cutoff := time.Now().Add(-time.Duration(beforeDays) * 24 * time.Hour)
	cutoff := mtime.Now().AddDays(-beforeDays)

	// 准备要匹配的类型集合
	typesToMatch := map[string]struct{}{}
	if len(opt.Type) == 0 {
		for k := range validTypes {
			typesToMatch[k] = struct{}{}
		}
	} else {
		for _, t := range opt.Type {
			tt := strings.ToLower(strings.TrimSpace(t))
			if _, ok := validTypes[tt]; ok {
				typesToMatch[tt] = struct{}{}
			}
		}
	}

	fileList, err := mfile.ListDir(l.Path, 0)
	if err != nil {
		l.Error("err:mlog.Clear|mfile.ListDir", l, opt)
		return
	}

	for _, v := range fileList {
		fileName := v.Name
		// 名称格式：名称-类型-YYYY-MM-DD.log
		// 判断是否包含  l.Name
		if !strings.HasPrefix(fileName, l.Name+"-") || !strings.HasSuffix(fileName, ".log") {
			continue
		}

		// 去掉前缀和后缀，得到 中间部分
		mid := strings.TrimSuffix(strings.TrimPrefix(fileName, l.Name+"-"), ".log")
		// 是否包含 typesToMatch 中的类型
		parts := strings.SplitN(mid, "-", 2)
		if len(parts) != 2 {
			continue
		}
		fileType := parts[0]
		fileDateStr := parts[1]
		if _, ok := typesToMatch[fileType]; !ok {
			continue
		}

		// 判断日期是否早于 cutoff
		fileDate, err := mtime.Parse(fileDateStr)
		if err != nil {
			l.Error("err:mlog.Clear|mtime.Parse", l, opt)
			continue
		}
		timeDiff := fileDate.UnixMilli() - cutoff.UnixMilli()
		// 当前file 小于 cutoff 则跳过
		if timeDiff >= 0 {
			continue
		}

		if v.IsFile {
			err = os.Remove(v.AbsPath)
			if err != nil {
				l.Error("err:mlog.Clear|删除日志文件失败", l, opt)
				continue
			}
			l.Info("err:mlog.Clear|日志文件已删除", v.AbsPath)
		}
	}
}
