package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"net/http"

	"github.com/sunnywalden/sync-data/controllers"
	"github.com/sunnywalden/sync-data/pkg/auth"
	"github.com/sunnywalden/sync-data/pkg/types"
)

// GetToken, return token for platform user
func GetToken(c *gin.Context) {

	res := types.Response{
		Code: -1,
		Msg:  "request authkey nil",
		Data: nil,
	}
	var status = http.StatusBadRequest

	//获取所有请求参数
	platUser := c.Query("platuser")
	authKey := c.Query("authkey")
	if platUser == "" || authKey == "" {
		status = http.StatusBadRequest
		res.Code = -1
		res.Msg = errors.ErrQueryParamsNil.Error()
		c.JSON(
			status,
			res,
		)
	}

	user, err := controllers.SearchPlatUser(platUser)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		status = http.StatusInternalServerError
	}

	token,err := auth.GenerateToken(user, authKey)
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		status = http.StatusInternalServerError
	}

	res.Code = 0
	res.Msg = ""
	res.Data = map[string]string{
		"token": token,
	}

	c.JSON(
		status,
		res,
		)
}
