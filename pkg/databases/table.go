package databases

import (
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/types"
)

func initUserTable(db *gorm.DB) (err error) {
	if !db.HasTable(&models.User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&types.User{}).Error; err != nil {
			log.Fatalf("Create user table error!%s", err)
			return err
		}
	}

	return nil
}
