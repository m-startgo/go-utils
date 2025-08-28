package mencrypt

import (
	"testing"
)

// go test -v -run TestMo7
func TestMo7(t *testing.T) {
	id := UUID()
	idc := UUID()

	id2 := TimeID()
	id2c := TimeID()

	if id == idc {
		t.Fatal("UUID 生成重复")
	}
	if id2 == id2c {
		t.Fatal("TimeID 生成重复")
	}
}
