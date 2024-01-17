package server

import (
	"fmt"
	"github.com/mysll/toolkit"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	stdlog "log"
	"mediahub/internal/conf"
	"mediahub/internal/db"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func initConfig(option *conf.Options) {
	if ok, err := toolkit.PathExists(option.DataPath); !ok {
		err = os.Mkdir(option.DataPath, os.ModePerm)
		if err != nil {
			log.Fatalf("create dir failed, error %s", err.Error())
		}
	}
	configPath := filepath.Join(option.DataPath, "config.json")
	if ok, _ := toolkit.PathExists(configPath); ok {
		conf.LoadConfig(configPath)
		log.Infof("load config from %s", configPath)
	} else {
		config := conf.DefaultConfig()
		err := config.Save(configPath)
		if err != nil {
			log.Fatalf("save config failed, error %s", err.Error())
		}
	}
	log.Infof("init config")
}

func initDb() {
	logLevel := logger.Silent
	newLogger := logger.New(
		stdlog.New(log.StandardLogger().Out, "\r\n", stdlog.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	database := conf.GetConfig().Database
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: database.TablePrefix,
		},
		Logger: newLogger,
	}

	if !(strings.HasSuffix(database.DBFile, ".db") && len(database.DBFile) > 3) {
		log.Fatalf("db name error.")
	}

	dB, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s?_journal=WAL&_vacuum=incremental",
		database.DBFile)), gormConfig)

	if err != nil {
		log.Fatalf("open database faield, error %s", err.Error())
	}

	db.InitDb(dB)
	log.Infof("init db")
}

func preload(options *conf.Options) {
	log.Infof("MediaHub version: %s", conf.AppVersion)
	initConfig(options)
	initDb()
}

func Start(option *conf.Options) {
	preload(option)
	serve()
}

func Close() {
	log.Infof("shutdown server...")
	shutdown()
	db.Close()
	log.Infof("server exit")
}
