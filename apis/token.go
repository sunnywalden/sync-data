package apis

import (
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sunnywalden/sync-data/pkg/auth"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/errs"
	"github.com/sunnywalden/sync-data/pkg/types"
)

// GetToken, return token for platform user
func GetToken(c *gin.Context) {

	var log *logrus.Logger
	configures := config.Conf
	log = logging.GetLogger(configures.Log.Level)

	res := types.Response{
		Code: -1,
		Msg:  "request authkey nil",
		Data: nil,
	}
	var status = http.StatusBadRequest

	//获取所有请求参数
	platUser := c.PostForm("platuser")
	authKey := c.PostForm("authkey")
	if platUser == "" || authKey == "" {
		status = http.StatusBadRequest
		res.Code = -1
		res.Msg = errs.ErrQueryParamsNil.Error()
		c.JSON(
			status,
			res,
		)
	}

	user, err := databases.SearchPlatUser(platUser)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		status = http.StatusInternalServerError
		c.JSON(
			status,
			res,
		)
	}

	token, err := auth.GenerateToken(user)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		status = http.StatusInternalServerError
	} else {

		res.Code = 0
		res.Msg = "success"
		res.Data = map[string]string{
			"token": token,
		}
		log.Infof("token returned!")
	}
		c.JSON(
			status,
			res,
		)

}
