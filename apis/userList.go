package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/sync"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log = logging.GetLogger()
)

// UserList, query all users
//func UserList(w http.ResponseWriter, r *http.Request) {
func UserList(c *gin.Context) {

	//ctx, cancel := context.WithCancel(context.Background())

	//defer cancel()

	res := types.Response{
		Code: 0,
		Msg: "",
		Data: nil,
	}
	var status = http.StatusOK

	configures := config.Conf

	users, err := sync.GetUser(c, configures)
	//users, err := sync.GetUser(ctx, configures)
	if err != nil {
		//log.Error("Getting all users err!%s", err)
		status = http.StatusInternalServerError
		res.Msg = err.Error()
		res.Code = -1
	} else {
		res.Data = users
	}

	//helper.ResponseWithJson(w, status, res)
	c.JSON(
		status,
		res,
		)

}
