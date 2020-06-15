package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/sync"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log *logrus.Logger
)

// UserList, query all users
func UserList(c *gin.Context) {

	res := types.Response{
		Code: 0,
		Msg: "",
		Data: nil,
	}
	var status = http.StatusOK

	configures := config.Conf

	log = logging.GetLogger(&configures.Log)

	users, err := sync.GetUser(c, configures)
	if err != nil {
		status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Code = -1
	} else {
		res.Data = users
	}

	c.JSON(
		status,
		res,
		)

}
