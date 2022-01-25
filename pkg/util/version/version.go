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
	if len(v1s) != 3 || len(v2s) != 3 {
		return false
	}

	for i := 0; i < 3; i++ {
		v1si, err := strconv.Atoi(v1s[i])
		if err != nil {
			return false
		}
		v2si, err := strconv.Atoi(v2s[i])
		if err != nil {
			return false
		}
		if v1si > v2si {
			return true
		}
	}
	return false
}
