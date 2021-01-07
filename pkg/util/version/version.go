package version

import (
	"strconv"
	"strings"
)

func IsNewerThan(v1, v2 string) bool {
	v1 = strings.Replace(v1, "v", "", -1)
	v2 = strings.Replace(v2, "v", "", -1)
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	for i := 0; i < 3; i++ {
		v1si, _ := strconv.Atoi(v1s[i])
		v2si, _ := strconv.Atoi(v2s[i])
		if v1si > v2si {
			return true
		}
	}
	return false
}
