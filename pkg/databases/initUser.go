package databases

import (
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/pkg/models"
)

// initUserTable, create user table if not exists
func initUserTable(db *gorm.DB) (err error) {
	if !db.HasTable(&models.User{}) {
		if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").CreateTable(&models.User{}).Error; err != nil {
			log.Fatalf("Create user table error!%s", err)
			return err
		}
	}

	return nil
}
