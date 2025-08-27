package mjson

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// 将任意 JSON-able 数据转换为 JSON 字节切片。
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
	err2 := json.Indent(&out, jsonByte, "", " ")
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

	err2 := json.Unmarshal(jsonByte, &resData)
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
