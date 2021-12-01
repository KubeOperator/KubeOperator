package lang

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/mozillazg/go-pinyin"
	"regexp"
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

func GetStringKey(val string) string {
	sha := sha1.New()
	sha.Write([]byte(val))
	result := hex.EncodeToString(sha.Sum([]byte("")))
	reg := regexp.MustCompile(`[0-9]`)
	return reg.ReplaceAllString(result, "")
}
