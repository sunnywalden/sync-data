package databases

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
)

var (
	log = logging.GetLogger()
)

func Conn(configures *config.MysqlConf) (db *gorm.DB,err error) {

	mysqlConf := configures
	mysqlHost := mysqlConf.Host + ":" + mysqlConf.Port
	mysqlDB := mysqlConf.DB
	mysqlUser := mysqlConf.User
	mysqlPassword := mysqlConf.Password

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlUser, mysqlPassword, mysqlHost, mysqlDB)
	log.Printf("Debug mysql connect string:%s\n", dbUrl)

	db, err = gorm.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalf("Connect mysql error!%s", err)
		return nil, err
	}
	return db, nil
}

