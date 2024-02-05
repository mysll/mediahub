package downloader

import (
	"errors"
	qbt "github.com/NullpointerW/go-qbittorrent-apiv2"
	log "github.com/sirupsen/logrus"
)

const (
	Downloading = "downloading"
	Completed   = "completed"
)

var (
	ErrNoClient = errors.New("qbittorrent client is nil")
)

type QBittorrent struct {
	client *qbt.Client
}

func NewQBittorrent(url string, account, password string) *QBittorrent {
	client, err := qbt.NewCli(url, account, password)
	if err != nil {
		log.Errorf("create qbittorrent client failed, %s", err.Error())
	}
	qb := &QBittorrent{
		client: client,
	}
	return qb
}

func (qb *QBittorrent) GetTorrent(hash string) (*qbt.TorrentProp, error) {
	if qb.client == nil {
		return nil, ErrNoClient
	}
	prop, err := qb.client.GetTorrentProperties(hash)
	if err != nil {
		return nil, err
	}
	return &prop, nil
}

func (qb *QBittorrent) GetFiles(hash string) ([]qbt.File, error) {
	if qb.client == nil {
		return nil, ErrNoClient
	}
	return qb.client.Files(hash)
}

func (qb *QBittorrent) AddTag(hash string, tag string) error {
	if qb.client == nil {
		return ErrNoClient
	}
	return nil
}

func (qb *QBittorrent) GetTorrentList(filter string) ([]qbt.Torrent, error) {
	if qb.client == nil {
		return nil, ErrNoClient
	}
	opts := make(qbt.Optional)
	if filter != "" {
		opts["filter"] = filter
	}
	return qb.client.TorrentList(opts)
}
