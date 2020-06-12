package databases

import (
	"github.com/sunnywalden/sync-data/config"
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/pkg/types"
)


var db *gorm.DB

// Init, init mysql oa database connection and create user table if not exists
func Init(configures *config.MysqlConf) () {
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
func UpdateMysql(update map[string]byte,where map[string]byte, configures *config.MysqlConf) (database *gorm.DB, err error) {

	Init(configures)

	for k,v := range update {

		database = db.Model(&types.User{}).Update(k,v)
	}
	return database, nil
}
