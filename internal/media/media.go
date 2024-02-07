package media

import (
	"path"
	"strings"
)

const (
	MediaUnknown = iota
	MediaVideo
	MediaAnim
)

const (
	MediaTypeUnknown = iota
	MediaTypeMovie
	MediaTypeTv
)

var (
	Ext = [...]string{".mp4", ".mkv", ".ts", ".iso",
		".rmvb", ".avi", ".mov", ".mpeg",
		".mpg", ".wmv", ".3gp", ".asf",
		".m4v", ".flv", ".m2ts", ".strm",
		".tp"}
)

func IsMediaFile(f string) bool {
	ext := strings.ToLower(path.Ext(f))
	for _, e := range Ext {
		if e == ext {
			return true
		}
	}
	return false
}

type Media struct {
	tmdb *Tmdb
}

func (m *Media) GetMediaInfo(title string, subtitle string) MetaInfo {
	var media MetaInfo
	if m.tmdb == nil {
		return nil
	}

	return media
}
