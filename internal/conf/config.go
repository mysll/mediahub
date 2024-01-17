package conf

import (
	"encoding/json"
	"github.com/mysll/toolkit"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var config *Config

func GetConfig() *Config {
	return config
}

type Database struct {
	Type        string `json:"type" env:"TYPE"`
	Host        string `json:"host" env:"HOST"`
	Port        int    `json:"port" env:"PORT"`
	User        string `json:"user" env:"USER"`
	Password    string `json:"password" env:"PASS"`
	Name        string `json:"name" env:"NAME"`
	DBFile      string `json:"db_file" env:"FILE"`
	TablePrefix string `json:"table_prefix" env:"TABLE_PREFIX"`
	SSLMode     string `json:"ssl_mode" env:"SSL_MODE"`
}

type Server struct {
	Address  string `json:"address" env:"ADDR"`
	HttpPort int    `json:"http_port" env:"HTTP_PORT"`
}

type Cors struct {
	AllowOrigins []string `json:"allow_origins" env:"ALLOW_ORIGINS"`
	AllowMethods []string `json:"allow_methods" env:"ALLOW_METHODS"`
	AllowHeaders []string `json:"allow_headers" env:"ALLOW_HEADERS"`
}

type Config struct {
	Database Database `json:"database"`
	Server   Server   `json:"server"`
	Cors     Cors     `json:"cors" envPrefix:"CORS_"`
}

func (c *Config) Load(f string) {
	data, err := toolkit.ReadFile(f)
	if err != nil {
		log.Fatalf("load config failed, %s", err.Error())
	}
	if err = json.Unmarshal(data, c); err != nil {
		log.Fatalf("load config failed, %s", err.Error())
	}
}

func (c *Config) Save(f string) error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(f, bytes, os.ModePerm)
}
func DefaultConfig() *Config {
	dbPath := filepath.Join(options.DataPath, "data.db")
	config = &Config{
		Server: Server{
			Address:  "0.0.0.0",
			HttpPort: 3005,
		},
		Database: Database{
			DBFile:      dbPath,
			TablePrefix: "mh_",
		},
		Cors: Cors{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		},
	}
	return config
}

func LoadConfig(f string) *Config {
	config = DefaultConfig()
	config.Load(f)
	return config
}
