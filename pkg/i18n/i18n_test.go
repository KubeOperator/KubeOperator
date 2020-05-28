package i18n

import (
	"testing"
)

func TestLang(t *testing.T) {

	m := make(map[string]string)
	m["name"] = "张三"
	println(Tr("invalid_test", m))

}
