package databases

import (
	"database/sql"
	"github.com/jinzhu/gorm"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/models"
)


var db *gorm.DB

// Init, init mysql oa database connection and create user table if not exists
func Init(configures *config.TomlConfig) (err error) {
	//var err error


	db, err = Conn(configures)
	if err != nil {
		return err
		//panic(err)
	}

	err = initUserTable(db)
	if err != nil {
		//panic(err)
		return err
	}

	err = InitPlatUserTable(db)
	if err != nil {
		return err
		//panic(err)
	}
	return nil
}

// UpdateMysql, update user table
func UpdateMysql(update map[string]byte,where map[string]byte, configures *config.TomlConfig) (database *gorm.DB, err error) {

	Init(configures)

	for k,v := range update {

		database = db.Model(&models.User{}).Update(k,v)
	}
	return database, nil
}


// UpdatePlat, platform user create
func UpdatePlat(user *models.PlatUser) (rows  *sql.Rows,err error) {
	configures := config.Conf
	db,err := Conn(configures)
	if err != nil {
		return nil, err
	}

	err = Init(configures)
	if err != nil {
		return nil, err
		//res.Msg = "database platform user create err"
	}
	log.Debugf("platform user table created!")

	rows, err = db.Create(user).Rows()
	return rows, err
}