package mlog

/*

myLog := mlog.New(Config{
	Path: "./logs",
	Name: "log",
})

myLog.Clear(ClearOpt{
	Type:   []string{"debug", "warn"},
	Before: 7,
})

*/

import (
	"os"
	"strings"

	"github.com/m-startgo/go-utils/mfile"
	"github.com/m-startgo/go-utils/mtime"
)

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
