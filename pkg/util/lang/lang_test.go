package lang

import (
	"fmt"
	"testing"
)

var val = "张三的mac"

func TestPinyin(t *testing.T) {
	fmt.Println(Pinyin(val))
}

