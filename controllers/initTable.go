package controllers

import (
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/pkg/models"
)

func InitPlatUserTable(db *gorm.DB) (err error) {
	if !db.HasTable(&models.User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.PlatUser{}).Error; err != nil {
			log.Errorf("Create platform user table error!%s", err)
			return err
		}
	}

	return nil
}
