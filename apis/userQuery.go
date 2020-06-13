package apis

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/sync"
	"github.com/sunnywalden/sync-data/pkg/types"
)


// searchUser, query user using params
func searchUser(ctx context.Context, attr string, searchStr string) (user *models.User, err error) {
	users, err := sync.GetUser(ctx, config.Conf)
	if err != nil {
		return nil, err
	}

	// 查询匹配的用户
	for _,user := range users {
		switch attr {
		case "nickName":
			if user.NickName == searchStr {
				log.Infof("User matched :%s", user.NickName)
				return &user, nil
			}
		case "loginId":
			if user.LoginId == searchStr {
				log.Infof("User matched %s:%s", user.LoginId, user.NickName)
				return &user, nil
			}
		default:
			if user.Name == searchStr {
				log.Infof("User matched %s:%s", user.Name, user.NickName)
				return &user, nil
			}
		}

	}

	return nil, errors.ErrUserNotExists

}


// User, query user matched
//func User(w http.ResponseWriter, r *http.Request) {
func User(c *gin.Context) {

	var (
		user *models.User
		err error
	)

	//ctx, cancel := context.WithCancel(context.Background())

	//defer cancel()

	res := types.Response{
		Code: 0,
		Msg: "",
		Data: nil,
	}
	var status = http.StatusOK

	//获取所有请求参数
	//query := r.URL.Query()

	nickName := c.DefaultQuery("nickname", "")
	userId := c.Query("id")
	name := c.Query("name")
	if nickName == "" && userId == "" && name == "" {
		status = http.StatusBadRequest
		res.Code = -1
		res.Msg = errors.ErrQueryParamsNil.Error()
		c.JSON(
			status,
			res,
			)
	}

	//nickName, ok := query["nickname"]
	//if ok{
	//	log.Printf("Debug user nickname %s", nickName[0])
	//	user, err = searchUser(ctx, "nickName", nickName[0])
	//}
	//userId, ok := query["id"]
	//if ok {
	//	log.Infof("Debug user loginId %s", userId[0])
	//	user, err = searchUser(ctx, "loginId", userId[0])
	//}

	//name, ok := query["name"]
	//if ok {
	//	log.Infof("Debug user lastname %s", name[0])
	//	user, err = searchUser(ctx, "lastName", name[0])
	//} else {
	//	status = http.StatusBadRequest
	//	res.Code = -1
	//	res.Msg = errors.ErrQueryParamsNil.Error()
	//}
	if nickName != "" {
		user, err = searchUser(c, "nickName", nickName)
	} else if userId != "" {
		user, err = searchUser(c, "loginId", userId)
	} else {
		user, err = searchUser(c, "lastName", name)
	}


	// 判断查询用户是否异常
	if err != nil {
		status = http.StatusInternalServerError
		res.Code = -1
		res.Msg = err.Error()
	} else {
		res.Data = user
	}

	//helper.ResponseWithJson(w, status, res)
	c.JSON(
		status,
		res,
		)
}
