package apis

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/controllers"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/types"
)

//func Register(w http.ResponseWriter, r *http.Request) {
func Register(c *gin.Context) {
	var user models.PlatUser
	res := types.Response{
		Code: http.StatusInternalServerError,
		Msg: "internal error",
		Data: nil,
	}
	var status = http.StatusInternalServerError

	platUser := c.Query("platuser")
	authKey := c.Query("authkey")
	user.UserName = platUser
	user.AuthKey = authKey

	//err := json.NewDecoder(r.Body).Decode(&user)
	//if err != nil || user.UserName == "" || user.AuthKey == "" {
	if user.UserName == "" || user.AuthKey == "" {
		res.Code = http.StatusBadRequest
		res.Msg = "bad params"
		c.JSON(
			status,
			res,
		)
		//helper.ResponseWithJson(
		//	w,
		//	http.StatusBadRequest,
		//	res,
		//)
		//return
	}

	configures := config.Conf

	db,err := databases.Conn(&configures.Mysql)
	if err != nil {
		res.Msg = "database connect err"
		//helper.ResponseWithJson(
		//	w,
		//	http.StatusInternalServerError,
		//	res,
		//)
		c.JSON(
			status,
			res,
		)
	}

	err = controllers.InitPlatUserTable(db)
	if err != nil {
		res.Msg = "database platform user create err"
		c.JSON(
			status,
			res,
		)
		//helper.ResponseWithJson(
		//	w,
		//	http.StatusInternalServerError,
		//	res,
		//)
	}

	rows, err := db.Model(&models.PlatUser{}).Create(user).Rows()
	if err != nil {
		res.Msg = "update database platform user err"
		//helper.ResponseWithJson(
		//	w,
		//	http.StatusInternalServerError,
		//	res,
		//)
		c.JSON(
			status,
			res,
		)
	}

	res.Msg = ""
	res.Code = http.StatusOK
	res.Data = rows
	//helper.ResponseWithJson(
	//	w,
	//	http.StatusOK,
	//	res,
	//)
	c.JSON(
		status,
		res,
		)
}
