package ssh

import "os"

func CheckIllegal(args ...string) bool {
	if args == nil {
		return false
	}
	for _, arg := range args {
		_, err := os.Stat(arg)
		if err != nil || os.IsNotExist(err) {
			return false
		}
	}
	return true
}
