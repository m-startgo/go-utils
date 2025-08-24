# go-utils

基于 go 语言的，常用的，通用的，开箱即用的，工具函数/辅助函数

包含很多的工具跟模块,让 go 代码编写更加顺手。

## 安装方式

```bash

go get -u github.com/m-startgo/go-utils@latest

```

## 安装指定版本

```bash
go get -u github.com/m-startgo/go-utils@v0.0.7
```

## 使用示例

```go
// 引入指定目录如 m_str
import "github.com/m-startgo/go-utils/m_str"

// 使用指定函数如 m_str.TplFormat
func main() {
	tpl := `
app.name = ${appName}
app.ip = ${appIP}
app.port = ${appPort}
`
	data := map[string]string{
		"appName": "my_ap123p",
		"appIP":   "0.0.0.0",
	}

	s := m_str.TplFormat(tpl, data)
	fmt.Println("s", s)
}
```

其余模块及函数请查看各目录下的注释以及 `_test.go` 测试文件。
