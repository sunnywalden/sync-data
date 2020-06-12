package apis

import (
	"context"
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/sync"
)

var (
	log = logging.GetLogger()
)

// UserList, query all users
func UserList(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	res := map[string]interface{}{
		"code": 0,
		"msg": "",
		"data": nil,
	}

	configures := config.Conf

	users, err := sync.GetUser(ctx, configures)
	if err != nil {
		//log.Error("Getting all users err!%s", err)
		res["msg"] = err.Error()
		res["code"] = -1
	} else {
		res["data"] = users
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonRes, err := json.Marshal(&res)
	if err != nil {
		//log.Error("json encode user data err!%s", err)
		res["code"] = -1
		res["msg"] = err.Error()
	}


	_, err = fmt.Fprintf(w, "%s", jsonRes)
	if err != nil {
		log.Error("writing data to api err!%s", err)
		panic(err)
	}

}
