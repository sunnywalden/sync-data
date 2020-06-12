package sync

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/auth"
	"github.com/sunnywalden/sync-data/pkg/cache"
	"github.com/sunnywalden/sync-data/pkg/consts"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log = logging.GetLogger()
)

type UserInfo types.UserInfo
type User types.User

// jsonDecode, decode json type data
func jsonDecode(bodyC []byte) (types.UserInfo, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonMap types.UserInfo

	err := json.Unmarshal(bodyC, &jsonMap)
	if err != nil {
		log.Errorf("json decode err!%s\n", err.Error())
	} else {
		return jsonMap, nil
	}

	return types.UserInfo{}, err
}

// bodyResolve, resolve api response data
func bodyResolve(resp *http.Response) (types.UserInfo, error) {
	status := resp.StatusCode
	if status == 200 {
		defer resp.Body.Close()
		bodyC, _ := ioutil.ReadAll(resp.Body)
		jsonMap, err := jsonDecode(bodyC)
		if err != nil {
			return types.UserInfo{}, err
		}
		return jsonMap, nil
	}

	return types.UserInfo{}, errors.New(resp.Status)
}

// GetUser, query all active users of oa
func GetUser(ctx context.Context, configures *config.TomlConfig) (users []types.User,err error) {

	// 优先查询缓存
	users, err = getUsersFromCache(ctx, &configures.Redis)
	if err == nil && len(users) != 0 {
			return users, nil
	}

	// 查询数据库
	users, err = getUsersFromDB(&configures.Mysql)
	if err == nil && len(users) != 0 {
		return users, nil
	}

	// 查询API
	token,err := auth.GetToken()
	if err != nil {
		panic(err)
	} else {
		log.Debugf("oa token: %s.\n", token)
	}

	users, err = getUserFromApi(token, &configures.OA)
	if err != nil {
		return nil, err
	}

	log.Debugf("Debug users in GetUser:%s\n", users[:5])

	// 刷新缓存
	err = setUsersCache(ctx, users, &configures.Redis)
	// 更新数据库
	err = storeUsersToDB(&configures.Mysql, users)

	return users, err
}

// getUserFromApi, get oa users from oa user query api
func getUserFromApi(token types.OaToken, configures *config.OA) (users []types.User,err error) {
	oaConf := configures
	oaUrl := oaConf.OaUrl
	oaUserApi := oaConf.UserApi
	activeUser := oaConf.UserActive

	log.Infof("Getting users from api!")
	queryUrl := fmt.Sprintf("%s%s?status=%d", oaUrl, oaUserApi, activeUser)
	log.Debugf("Debug oa user api url:%s\n", queryUrl)

	req, err := http.NewRequest("GET", queryUrl, nil)

	if err != nil {
		log.Errorf("http request error!%s\n", err)
		return nil, err
	}
	req.Header.Set("authorization", "Bearer " + string(token))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errorf("query oa user api error!%s\n", err)
		return nil, err
	}

	res, err := bodyResolve(resp)
	if err != nil {
		return nil, err
	}

	code := res.Code
	if code != consts.OAOkCode {
		return nil, errors.New(res.Message)
	}

	data := res.Data
	return data, nil
}

// getUsersFromCache, get all users from cache
func getUsersFromCache(ctx context.Context, configures *config.RedisConf) (users []types.User,err error) {

	userKey := consts.UserInfoKey

	log.Infof("Getting users from cache!")
	client, err := cache.GetClient(ctx, configures)
	if err != nil {
		log.Errorf("Get cache client error!%s\n", err)
		return nil, err
	}

	val, err := client.Get(ctx, userKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("Get users from cache error!%s\n", err)
		}
		return nil, err
	}

	//logging.Printf("Debug users before decode:%s\n", val)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		log.Errorf("Decode users from cache error!%s\n", err)
		return nil, err
	}

	return users, nil

}

// getUsersFromDB, get all users from mysql user table
func getUsersFromDB(configures *config.MysqlConf) (users []types.User,err error) {
	db, err := databases.Conn(configures)
	if err != nil {
		return nil, err
	}

	databases.Init(configures)

	usersRows, err := db.Model(&types.User{}).Rows()
	if err != nil {
		log.Errorf("Query user table error!%s\n", err)
		return nil, err
	}

	log.Debugf("%s\n", usersRows)

	return nil,err
}

// storeUsersToDB, insert all users to mysql user table
func storeUsersToDB(configures *config.MysqlConf, users []types.User) (err error) {
	db, err := databases.Conn(configures)
	if err != nil {
		return err
	}

	databases.Init(configures)

	var userDb *gorm.DB

	for id, userInfo := range users {
		userInfo.Id = id
		userDb = db.Create(userInfo)
	}

	if userDb != nil {
		log.Debugf("%d\n", userDb.RowsAffected)
	}

	return nil
}

// setUsersCache, store all users list to cache
func setUsersCache(ctx context.Context, users []types.User, configures *config.RedisConf) (err error) {

	userKey := consts.UserInfoKey
	expireDay := consts.UserCacheExpireDay
	var Days = time.Duration(expireDay)

	log.Info("Store users to cache!")
	client, err := cache.GetClient(ctx, configures)
	if err != nil {
		log.Errorf("Get redis client error!%s\n", err)
		return err
	}

	err = client.Set(ctx, "sayhello", "hello world", time.Hour * 24 * Days).Err()
	if err != nil {
		log.Errorf("Redis set test error!%s\n", err)
		return err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	usersByte, err := json.Marshal(users)
	if err != nil {
		log.Errorf("Json encode users error!%s\n", err)
	}

	err = client.Set(ctx, userKey, usersByte, time.Hour * 24 * Days).Err()
	if err != nil {
		log.Errorf("Store users to cache error!%s\n", err)
		return  err
	}
	return nil
}