package mjson

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ToByte 将任意数据序列化为 JSON bytes。
// 功能：将 data 使用 json.Marshal 转为 []byte。
// 参数：data any — 要序列化的数据。
// 返回：[]byte, error — 序列化后的字节数组与可能的错误。
// 错误格式：err:mjson.ToByte|marshal|<底层错误>
// 使用示例：
// b, err := mjson.ToByte(map[string]any{"a":1})
// if err != nil { /* 处理错误 */ }
func ToByte(data any) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("err:mjson.ToByte|nil|data is nil")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("err:mjson.ToByte|marshal|%w", err)
	}
	return b, nil
}

// ToStr 将任意数据序列化为 JSON 字符串；出错时返回空 JSON (`{}`) 并吞掉错误。
// 功能：将 data 转为 JSON 字符串。
// 参数：data any。
// 返回：string — JSON 字符串，发生错误时返回 "{}"。
// 使用示例：
// s := mjson.ToStr([]int{1,2,3})
func ToStr(data any) string {
	b, err := ToByte(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// IndentByteToStr 将 JSON bytes 格式化为漂亮的 JSON 字符串（使用 json.Indent）。
// 功能：对传入的 JSON bytes 使用 json.Indent 缩进格式化并返回字符串。
// 参数：b []byte — 原始 JSON bytes。
// 返回：string, error — 格式化后的 JSON 字符串与可能的错误。
// 错误格式：err:mjson.IndentByteToStr|indent|<底层错误>
// 使用示例：
// s, err := mjson.IndentByteToStr([]byte(`{"a":1}`))
// IndentByteToStr 将 JSON bytes 格式化为漂亮的 JSON 字符串（使用 json.Indent）。
// 出错时返回空 JSON (`{}`) 并吞掉错误。
func IndentByteToStr(b []byte) string {
	if len(b) == 0 {
		return "{}"
	}
	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		return "{}"
	}
	return out.String()
}

// IndentAnyToString 将任意数据先序列化，再使用 json.Indent 返回格式化后的 JSON 字符串。
// 功能：对 data 进行 Marshal，然后使用 json.Indent 美化输出。
// 参数：data any。
// 返回：string, error。
// 使用示例：
// s, err := mjson.IndentAnyToString(struct{A int}{A:1})
// IndentAnyToString 将任意数据先序列化，再使用 json.Indent 返回格式化后的 JSON 字符串。
// 出错时返回空 JSON (`{}`) 并吞掉错误。
func IndentAnyToString(data any) string {
	// 如果 data 已经是 []byte 或 string，先转换为 []byte 再处理
	switch v := data.(type) {
	case []byte:
		return IndentByteToStr(v)
	case string:
		return IndentByteToStr([]byte(v))
	}

	b, err := ToByte(data)
	if err != nil {
		return "{}"
	}
	return IndentByteToStr(b)
}

// ToMap 将任意 JSON-able 数据转换为 map[string]any。
// 功能：支持传入 struct/map/string/[]byte 等类型，将其转换为 map[string]any。
// 参数：data any。
// 返回：map[string]any, error — 转换结果或错误。
// 常见错误：当输入无法解析为对象时返回错误，格式：err:mjson.ToMap|unmarshal|<底层错误>
// 使用示例：
// m, err := mjson.ToMap(`{"a":1}`)
func ToMap(val any) (resData map[string]any, resErr error) {
	resData = map[string]any{}
	resErr = nil

	jsonByte, err := ToByte(val)
	if err != nil {
		resErr = err
		return
	}
	err2 := json.Unmarshal(jsonByte, &resData)
	if err2 != nil {
		resData = map[string]any{}
		return
	}
	return
}

// PrintAny 将任意数据美化为 JSON 并打印到标准输出。
// 功能：对 data 使用 IndentAnyToString，并打印结果。
// 参数：data any。
// 错误处理：函数不返回 error，如遇到错误会直接打印错误信息。
// 使用示例：
// mjson.PrintAny(map[string]any{"a":1})
func PrintAny(data any) {
	s := IndentAnyToString(data)
	fmt.Println(s)
}

// PrintByte 将 JSON bytes 美化并打印到标准输出。
// 功能：对 JSON bytes 使用 json.Indent 并打印结果。
// 参数：b []byte。
// 错误处理：函数不返回 error，如遇到错误会直接打印错误信息。
// 使用示例：
// mjson.PrintByte([]byte(`{"a":1}`))
func PrintByte(b []byte) {
	s := IndentByteToStr(b)
	fmt.Println(s)
}
