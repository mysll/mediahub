package media

import (
	"fmt"
	tmdb "github.com/cyruzin/golang-tmdb"
	"github.com/hashicorp/golang-lru/arc/v2"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

type Tmdb struct {
	client      *tmdb.Client
	options     map[string]string
	movieCache  *arc.ARCCache[int, *tmdb.MovieDetails]
	tvCache     *arc.ARCCache[int, *tmdb.TVDetails]
	searchCache *arc.ARCCache[string, *tmdb.SearchMulti]
}

func NewTmdb(apiKey string, language string, proxyUrl string) *Tmdb {
	tmdbClient, err := tmdb.Init(apiKey)
	if err != nil {
		log.Fatalf("create tmdb failed, err %s", err.Error())
	}
	options := make(map[string]string)
	if language == "" {
		language = "zh"
	}
	options["language"] = language
	var proxy func(*http.Request) (*url.URL, error)
	if proxyUrl != "" {
		proxyUrl_, err := url.Parse(proxyUrl)
		if err != nil {
			log.Errorf("parse proxy url failed, %s", err.Error())
		} else {
			proxy = http.ProxyURL(proxyUrl_)
		}
	}
	tmdbClient.SetClientAutoRetry()
	tmdbClient.SetClientConfig(http.Client{
		Timeout: time.Second * 5,
		Transport: &http.Transport{
			Proxy:           proxy,
			MaxIdleConns:    10,
			IdleConnTimeout: 15 * time.Second,
		},
	})
	detailCache, err := arc.NewARC[int, *tmdb.MovieDetails](512)
	if err != nil {
		log.Fatal(err.Error())
	}
	searchCache, err := arc.NewARC[string, *tmdb.SearchMulti](512)
	if err != nil {
		log.Fatal(err.Error())
	}
	tvCache, err := arc.NewARC[int, *tmdb.TVDetails](512)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &Tmdb{
		client:      tmdbClient,
		options:     options,
		movieCache:  detailCache,
		searchCache: searchCache,
		tvCache:     tvCache,
	}
}

func (t *Tmdb) QueryByName(name string) (err error) {
	search, ok := t.searchCache.Get(name)
	if !ok {
		search, err = t.client.GetSearchMulti(name, t.options)
		if err != nil {
			return err
		}
		t.searchCache.Add(name, search)
	}

	for _, result := range search.Results {
		if result.MediaType == "movie" {
			fmt.Println(result.MediaType, result.Title, result.ReleaseDate)
		} else if result.MediaType == "tv" {
			fmt.Println(result.MediaType, result.Name, result.FirstAirDate)
		}
	}
	return nil
}

func (t *Tmdb) GetMovieDetail(id int) (err error) {
	detail, ok := t.movieCache.Get(id)
	if !ok {
		detail, err = t.client.GetMovieDetails(id, t.options)
		if err != nil {
			return err
		}
		t.movieCache.Add(id, detail)
	}

	log.Infof("movie:%s overview:%s", detail.Title, detail.Overview)
	return nil
}

func (t *Tmdb) GetTvDetail(id int) (err error) {
	detail, ok := t.tvCache.Get(id)
	if !ok {
		detail, err = t.client.GetTVDetails(id, t.options)
		if err != nil {
			return err
		}
		t.tvCache.Add(id, detail)
	}

	log.Infof("tv:%s overview:%s", detail.Name, detail.Overview)
	return nil
}
