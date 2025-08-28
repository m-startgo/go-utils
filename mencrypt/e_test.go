package mencrypt

import (
	"fmt"
	"testing"
)

// go test -v -run TestMo7
func TestMo7(t *testing.T) {
	id := UUID()

	id2 := TimeID()

	fmt.Println("id:", id)
	fmt.Println("id2:", id2)
}
