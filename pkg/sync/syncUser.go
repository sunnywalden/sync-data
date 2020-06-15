package sync

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/auth"
	"github.com/sunnywalden/sync-data/pkg/cache"
	"github.com/sunnywalden/sync-data/pkg/consts"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log *logrus.Logger
)

type UserInfo types.UserInfo
type User models.User

// jsonDecode, decode json type data
func jsonDecode(bodyC []byte) (types.UserInfo, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonMap types.UserInfo

	log.Debugf("debug data before decode:%s", bodyC)
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
func GetUser(ctx context.Context, configures *config.TomlConfig) (users []models.User,err error) {

	log = logging.GetLogger(&configures.Log)

	// 优先查询缓存
	users, err = getUsersFromCache(ctx, configures)
	if err == nil && len(users) != 0 {
			return users, nil
	}

	// 查询数据库
	users, err = getUsersFromDB(configures)
	if err == nil && len(users) != 0 {
		return users, nil
	}

	// 查询API
	token,err := auth.GetToken()
	if err != nil {
		panic(err)
	} else {
		log.Debugf("oa token: %s", token)
	}

	users, err = getUserFromApi(token, &configures.OA)
	if err != nil {
		return nil, err
	}

	log.Debugf("Debug users in GetUser:%s", users[:5])

	// 刷新缓存
	err = setUsersCache(ctx, users, configures)
	// 更新数据库
	err = storeUsersToDB(configures, users)

	return users, err
}

// getUserFromApi, get oa users from oa user query api
func getUserFromApi(token types.OaToken, configures *config.OA) (users []models.User,err error) {
	oaConf := configures
	oaUrl := oaConf.OaUrl
	oaUserApi := oaConf.UserApi
	activeUser := oaConf.UserActive

	log.Infof("Getting users from api!")
	queryUrl := fmt.Sprintf("%s%s?status=%d", oaUrl, oaUserApi, activeUser)
	log.Debugf("Debug oa user api url:%s", queryUrl)

	req, err := http.NewRequest("GET", queryUrl, nil)

	if err != nil {
		log.Errorf("http request error!%s", err)
		return nil, err
	}
	req.Header.Set("authorization", "Bearer " + string(token))

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Errorf("query oa user api error!%s", err)
		return nil, err
	}

	log.Debugf("debug response from api:%d", resp.StatusCode)
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
func getUsersFromCache(ctx context.Context, configures *config.TomlConfig) (users []models.User,err error) {

	userKey := consts.UserInfoKey

	log.Infof("Getting users from cache!")
	client, err := cache.GetClient(ctx, configures)
	if err != nil {
		log.Errorf("Get cache client error!%s", err)
		return nil, err
	}

	val, err := client.Get(ctx, userKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Errorf("Get users from cache error!%s", err)
		}
		return nil, err
	}

	//logging.Printf("Debug users before decode:%s\n", val)

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal([]byte(val), &users)
	if err != nil {
		log.Errorf("Decode users from cache error!%s", err)
		return nil, err
	}

	return users, nil

}

// getUsersFromDB, get all users from mysql user table
func getUsersFromDB(configures *config.TomlConfig) (users []models.User,err error) {
	db, err := databases.Conn(configures)
	if err != nil {
		return nil, err
	}

	databases.Init(configures)

	usersRows, err := db.Model(&models.User{}).Rows()
	if err != nil {
		log.Errorf("Query user table error!%s", err)
		return nil, err
	}

	log.Debugf("%s", usersRows)

	return nil,err
}

// storeUsersToDB, insert all users to mysql user table
func storeUsersToDB(configures *config.TomlConfig, users []models.User) (err error) {
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
		log.Debugf("%d", userDb.RowsAffected)
	}

	return nil
}

// setUsersCache, store all users list to cache
func setUsersCache(ctx context.Context, users []models.User, configures *config.TomlConfig) (err error) {

	userKey := consts.UserInfoKey
	expireDay := consts.UserCacheExpireDay
	var Days = time.Duration(expireDay)

	log.Info("Store users to cache!")
	client, err := cache.GetClient(ctx, configures)
	if err != nil {
		log.Errorf("Get redis client error!%s", err)
		return err
	}

	err = client.Set(ctx, "sayhello", "hello world", time.Hour * 24 * Days).Err()
	if err != nil {
		log.Errorf("Redis set test error!%s", err)
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
