package utils

import "github.com/dlclark/regexp2"

func MatchString(re *regexp2.Regexp, text string) bool {
	ok, _ := re.MatchString(text)
	return ok
}

func ReplaceString(re *regexp2.Regexp, text string, replace string) string {
	if t, err := re.Replace(text, replace, -1, -1); err == nil {
		return t
	}
	return text
}
