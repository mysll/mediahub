package media

import (
	"testing"
)

func TestTmdb(t *testing.T) {
	tmdb := NewTmdb("7fcd52ad8f1f1801e55def029f0b5f09", "zh", "http://127.0.0.1:7890")
	err := tmdb.QueryByName("冰血暴")
	if err != nil {
		t.Fatal(err.Error())
	}
}
