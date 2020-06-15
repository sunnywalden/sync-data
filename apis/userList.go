package apis

import (
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/sync"
	"github.com/sunnywalden/sync-data/pkg/types"
)

// UserList, query all users
func UserList(c *gin.Context) {

	var log *logrus.Logger

	configures := config.Conf

	log = logging.GetLogger(configures.Log.Level)

	res := types.Response{
		Code: 0,
		Msg: "",
		Data: nil,
	}
	var status = http.StatusOK

	users, err := sync.GetUser(c, configures)
	if err != nil {
		status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Code = -1
	} else {
		res.Data = users
		res.Msg = "success"
		log.Infof("all users query finished!")
	}

	c.JSON(
		status,
		res,
		)

}
