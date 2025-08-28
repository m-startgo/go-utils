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
