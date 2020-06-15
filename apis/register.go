package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/types"
)

func Register(c *gin.Context) {

	var log *logrus.Logger
	configures := config.Conf
	log = logging.GetLogger(configures.Log.Level)

	var user models.PlatUser
	res := types.Response{
		Code: http.StatusInternalServerError,
		Msg: "internal error",
		Data: nil,
	}
	var status = http.StatusInternalServerError

	platUser := c.PostForm("platuser")
	authKey := c.PostForm("authkey")
	log.Debugf(platUser, authKey)

	if platUser == "" || authKey == "" {
		res.Code = http.StatusBadRequest
		res.Msg = "bad params"
		c.JSON(
			status,
			res,
		)
	}

	user.UserName = platUser
	user.AuthKey = authKey
	log.Debugf("Debug user info:%s", user)


	rows, err := databases.UpdatePlat(&user)
	if err != nil {
		res.Msg = "update database platform user err"
		c.JSON(
			status,
			res,
		)
	} else {
		res.Msg = ""
		res.Code = http.StatusOK
		res.Data = rows
		c.JSON(
			status,
			res,
		)
	}
}
