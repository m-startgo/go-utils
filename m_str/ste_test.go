package m_str

import (
	"fmt"
	"testing"
)

func Test_ToStr(t *testing.T) {
	// 测试 []rune 转字符串
	a := []rune("mo7欢迎你")
	aStr := ToStr(a)
	fmt.Println("a", aStr)
	t.Log(aStr)

	// 测试 []byte 转字符串
	b := []byte("mo7欢迎你")
	bStr := ToStr(b)
	fmt.Println("b", bStr)

	// 测试 float64 转字符串
	c := 10.97
	cStr := ToStr(c)
	fmt.Println("c", cStr)

	// 测试 int 转字符串
	d := 100
	dStr := ToStr(d)
	fmt.Println("d", dStr)

	// 测试 string 转字符串
	e := "hello"
	eStr := ToStr(e)
	fmt.Println("e", eStr)

	// 测试 bool 转字符串
	f := true
	fStr := ToStr(f)
	fmt.Println("f", fStr)

	// 测试 map 转字符串
	g := map[string]int{"one": 1, "two": 2}
	gStr := ToStr(g)
	fmt.Println("g", gStr)

	// 测试 struct 转字符串
	type Person struct {
		Name string
		Age  int
	}
	h := Person{Name: "Alice", Age: 30}
	hStr := ToStr(h)
	fmt.Println("h", hStr)

	// 测试 nil 转字符串
	var i any = nil
	iStr := ToStr(i)
	fmt.Println("i", iStr)

	// 测试 数组 转字符串
	j := [3]int{1, 2, 3}
	jStr := ToStr(j)
	fmt.Println("j", jStr)
}
