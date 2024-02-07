package utils

import "path"

func GetFileName(f string) string {
	_, file := path.Split(f)
	n := len(file)
	for i := 0; i < n; i++ {
		if file[i] == '.' {
			return file[:i]
		}
	}
	return file
}
