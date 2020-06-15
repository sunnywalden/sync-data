package databases

import (
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/models"
)


var db *gorm.DB

// Init, init mysql oa database connection and create user table if not exists
func Init(configures *config.TomlConfig) () {
	var err error


	db, err = Conn(configures)
	if err != nil {
		panic(err)
	}

	err = initUserTable(db)
	if err != nil {
		panic(err)
	}
}

// UpdateMysql, update user table
func UpdateMysql(update map[string]byte,where map[string]byte, configures *config.TomlConfig) (database *gorm.DB, err error) {

	Init(configures)

	for k,v := range update {

		database = db.Model(&models.User{}).Update(k,v)
	}
	return database, nil
}
