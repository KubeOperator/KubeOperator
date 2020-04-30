package util

import "fmt"

func SetFileContent(file, pattern, content string) string {
	return fmt.Sprintf("grep -Pq '%s' %s && sed -i 's;%s;%s;g' %s|| echo '%s' >> %s",
		pattern, file,
		pattern, content, file,
		content, file)
}
