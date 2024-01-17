package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"mediahub/internal/model"
)

var db *gorm.DB

func InitDb(d *gorm.DB) {
	db = d
	err := db.AutoMigrate(new(model.User))
	if err != nil {
		log.Fatalf("init db failed, error %s", err.Error())
	}
}

func GetDb() *gorm.DB {
	return db
}

func Close() {
	log.Info("closing db")
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf("failed to get db: %s", err.Error())
		return
	}
	err = sqlDB.Close()
	if err != nil {
		log.Errorf("failed to close db: %s", err.Error())
		return
	}
}
