package m_str

import (
	"testing"
)

func Test_ToStr(t *testing.T) {
	// 测试 []rune 转字符串
	a := []rune("mo7欢迎你")
	aStr := ToStr(a)
	t.Log(aStr)

	// 测试 []byte 转字符串
	b := []byte("mo7欢迎你")
	bStr := ToStr(b)
	t.Log(bStr)

	// 测试 float64 转字符串
	c := 10.97
	cStr := ToStr(c)
	t.Log(cStr)

	// 测试 int 转字符串
	d := 100
	dStr := ToStr(d)
	t.Log(dStr)

	// 测试 string 转字符串
	e := "hello"
	eStr := ToStr(e)
	t.Log(eStr)

	// 测试 bool 转字符串
	f := true
	fStr := ToStr(f)
	t.Log(fStr)

	// 测试 map 转字符串
	g := map[string]int{"one": 1, "two": 2}
	gStr := ToStr(g)
	t.Log(gStr)

	// 测试 struct 转字符串
	type Person struct {
		Name string
		Age  int
	}
	h := Person{Name: "Alice", Age: 30}
	hStr := ToStr(h)
	t.Log(hStr)

	// 测试 nil 转字符串
	var i any = nil
	iStr := ToStr(i)
	t.Log(iStr)

	// 测试 数组 转字符串
	j := [3]int{1, 2, 3}
	jStr := ToStr(j)
	t.Log(jStr)
}
