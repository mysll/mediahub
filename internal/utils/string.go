package utils

import (
	cn "github.com/aiialzy/chinese-number"
	"strconv"
	"strings"
)

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

func IsNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func CnToNumber(cs string, def int64) int64 {
	if v, err := strconv.ParseInt(cs, 10, 64); err == nil {
		return v
	}
	val, err := cn.Parse(cs)
	if err != nil {
		return def
	}
	return val
}
