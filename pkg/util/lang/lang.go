package lang

import (
	"github.com/mozillazg/go-pinyin"
	"unicode"
)

func Pinyin(hans string) string {
	a := pinyin.NewArgs()
	result := ""
	for _, v := range hans {
		if unicode.Is(unicode.Han, v) {
			p := pinyin.Slug(string(v), a)
			result += p
		} else {
			result += string(v)
		}
	}
	return result
}

func CountChinese(val string) int {
	count := 0
	for _, v := range val {
		if unicode.Is(unicode.Han, v) {
			count++
		}
	}
	return count
}
