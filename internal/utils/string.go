package utils

import "strings"

func IsChinese(s string) bool {
	for _, r := range s {
		if r == ' ' {
			continue
		}
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
}

func InList(s string, list []string, upper bool) bool {
	if upper {
		s = strings.ToUpper(s)
	}
	for _, ss := range list {
		if s == ss {
			return true
		}
	}
	return false
}
