package controllers

import (
	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/models"
)

var (
	log = logging.GetLogger()
)




// searchUser, query user using params
func SearchPlatUser(platUser string) (user *models.PlatUser, err error) {

	configures := config.Conf

	db, err := databases.Conn(&configures.Mysql)
	if err != nil {
		return nil, err
	}

	// 查询匹配的用户
	//var user models.PlatUser
	//queryDb := db.Model(&models.PlatUser{}).Select("username", map[string]string{"username": platUser})
	db.Where(models.PlatUser{UserName: platUser}).First(&user)

	if user != nil {
		return user, nil
	}
	return nil, errors.ErrUserNotExists

}