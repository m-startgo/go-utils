package mjson

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"log"

	jsoniter "github.com/json-iterator/go"
)

// Marshal 使用 jsoniter.ConfigCompatibleWithStandardLibrary 对 v 进行序列化。
// 保持与标准库兼容的语义，同时利用 jsoniter 的实现。
func Marshal(v any) ([]byte, error) {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(v)
}

// Unmarshal 使用 jsoniter.ConfigCompatibleWithStandardLibrary 对 data 进行反序列化到 v。
func Unmarshal(data []byte, v any) error {
	return jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal(data, v)
}

// 将任意 JSON-able 数据转换为 JSON 字节切片。
func ToByte(data any) ([]byte, error) {
	if data == nil {
		return nil, fmt.Errorf("err:mjson.ToByte|nil|data is nil")
	}
	b, err := Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("err:mjson.ToByte|marshal|%w", err)
	}
	return b, nil
}

// 将任意 JSON-able 数据转换为 JSON 字符串。
// 错误时返回 "{}"。
func ToStr(data any) string {
	b, err := ToByte(data)
	if err != nil {
		return "{}"
	}
	return string(b)
}

// 将任意 JSON-able 数据转换为缩进格式的 JSON 字符串。
// 错误时返回 "{}"。
func IndentJson(data any) string {
	jsonByte, err := ToByte(data)
	if err != nil {
		return "{}"
	}
	var out bytes.Buffer
	err2 := stdjson.Indent(&out, jsonByte, "", " ")
	if err2 != nil {
		return "{}"
	}
	return out.String()
}

// 将任意 JSON-able 数据转换为 map[string]any。
func ToMap(val any) (resData map[string]any, resErr error) {
	resData = map[string]any{}
	resErr = nil

	jsonByte, err := ToByte(val)
	if err != nil {
		resErr = err
		return
	}

	err2 := Unmarshal(jsonByte, &resData)
	if err2 != nil {
		resErr = err2
		return
	}
	return
}

// 将任意 JSON-able 数据转换为 map[string]string。
func ToMapStr(val any) (resData map[string]string, resErr error) {
	resData = map[string]string{}
	resErr = nil

	jsonByte, err := ToByte(val)
	if err != nil {
		resErr = err
		return
	}

	err2 := Unmarshal(jsonByte, &resData)
	if err2 != nil {
		resErr = err2
		return
	}
	return
}

// 打印任意 JSON-able 数据的缩进格式的 JSON 字符串，并返回该字符串。
func PrintAny(data any) string {
	s := IndentJson(data)
	log.Println(s)
	return s
}
