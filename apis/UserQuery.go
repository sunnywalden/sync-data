package apis

import (
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/errors"
	"github.com/sunnywalden/sync-data/pkg/sync"
	"github.com/sunnywalden/sync-data/pkg/types"
)


// searchUser, query user using params
func searchUser(ctx context.Context, attr string, searchStr string) (user *types.User, err error) {
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
func User(w http.ResponseWriter, r *http.Request) {

	var (
		user *types.User
		err error
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	res := map[string]interface{}{
		"code": -1,
		"msg": "",
		"data": nil,
	}

	//获取所有请求参数
	query := r.URL.Query()

	nickName, ok := query["nickname"]
	if ok{
		log.Printf("Debug user nickname %s", nickName[0])
		user, err = searchUser(ctx, "nickName", nickName[0])
	}
	userId, ok := query["id"]
	if ok {
		log.Infof("Debug user loginId %s", userId[0])
		user, err = searchUser(ctx, "loginId", userId[0])
	}

	name, ok := query["name"]
	if ok {
		log.Infof("Debug user lastname %s", name[0])
		user, err = searchUser(ctx, "lastName", name[0])
	} else {
			res["msg"] = errors.ErrQueryParamsNil
	}


	// 判断查询用户是否异常
	if err != nil {
			res["msg"] = err.Error()
	} else {
		res["code"] = 0
		res["data"] = user
	}

	// 结果转换为json格式
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonRes, err := json.Marshal(&res)
	if err != nil {
		res["msg"] = err.Error()
	}


	_, err = fmt.Fprintf(w, "%s", jsonRes)
	if err != nil {
		log.Error("writing data to api err!%s", err)
		panic(err)
	}
}
