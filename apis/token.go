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
//func GetToken(w http.ResponseWriter, r *http.Request) {
func GetToken(c *gin.Context) {

	res := types.Response{
		Code: -1,
		Msg:  "request authkey nil",
		Data: nil,
	}
	var status = http.StatusBadRequest

	//获取所有请求参数
	//query := r.URL.Query()
	//
	//platUser, ok := query["platuser"]
	//if !ok{
	//	res.Msg = "request platuser nil"
	//	helper.ResponseWithJson(
	//		w,
	//		status,
	//		res,
	//		)
	//}
	//log.Printf("Debug platuser %s", platUser[0])
	//
	//authKey, ok := query["authkey"]
	//if !ok{
	//	res.Msg = "request authkey nil"
	//	helper.ResponseWithJson(
	//		w,
	//		status,
	//		res,
	//	)
	//}
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
	//user, err := controllers.SearchPlatUser(platUser[0])
	if err != nil {
		res.Msg = err.Error()
		res.Code = -1
		status = http.StatusInternalServerError
	}

	token,err := auth.GenerateToken(user, authKey)
	//token,err := auth.GenerateToken(user, authKey[0])
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
	//helper.ResponseWithJson(
	//	w,
	//	status,
	//	res,
	//	)
}
