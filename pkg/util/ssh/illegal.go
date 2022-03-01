package ssh

import "strings"

func CheckIllegal(args string) bool {
	if strings.Contains(args, "&") || strings.Contains(args, "|") || strings.Contains(args, ";") ||
		strings.Contains(args, "$") || strings.Contains(args, "'") || strings.Contains(args, "`") ||
		strings.Contains(args, "(") || strings.Contains(args, ")") || strings.Contains(args, "\"") {
		return true
	}
	return false
}
