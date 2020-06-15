package databases

import (
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/pkg/errs"
	"github.com/sunnywalden/sync-data/pkg/logging"

	//"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/models"
)

var (
	log *logrus.Logger
	//log = logging.GetLogger()
)




// searchUser, query user using params
func SearchPlatUser(platUserName string) (user *models.PlatUser, err error) {

	var platUser models.PlatUser

	configures := config.Conf

	log = logging.GetLogger(configures.Log.Level)

	db, err := Conn(configures)
	if err != nil {
		return nil, err
	}

	// 查询匹配的用户
	db.Where(models.PlatUser{UserName: platUserName}).First(&platUser)

	if platUser.UserName != "" {
		return &platUser, nil
	}
	return nil, errs.ErrUserNotExists

}