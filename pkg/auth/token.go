package auth

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunnywalden/sync-data/pkg/logging"

	//"github.com/sunnywalden/sync-data/pkg/logging"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/consts"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log *logrus.Logger
	//log = logging.GetLogger()
	//log = config.Logger
	//log = logging.GetLogger(config.Conf.Log.Level)
)

// jsonDecode, json decode
func jsonDecode(bodyC []byte) (types.OAToken, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonMap types.OAToken

	err := json.Unmarshal(bodyC, &jsonMap)
	if err != nil {
		log.Fatalf("json decode err!%s", err.Error())
		return types.OAToken{}, err
	} else {
		return jsonMap, nil
	}
}

// bodyResolve, resolve response body to go struct
func bodyResolve(resp *http.Response) (types.OAToken, error) {
	status := resp.StatusCode
	if status == 200 {
		defer resp.Body.Close()
		bodyC, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Debug api response:%s", string(bodyC))
		jsonMap, err := jsonDecode(bodyC)
		if err != nil {
			log.Fatalf("Read api response error!%s", err)
			return types.OAToken{}, err
		}
		return jsonMap, nil
	}

	return types.OAToken{}, errors.New(resp.Status)
}

// get token from oa api
func GetToken() (token types.OaToken, err error) {
	configures := config.Conf
	log = logging.GetLogger(configures.Log.Level)
	oaConf := configures.OA
	oaUrl := oaConf.OaUrl
	oaTokenApi := oaConf.TokenApi
	oaClient := oaConf.OaClient
	oaSecret := oaConf.OaSecret

		queryUrl := fmt.Sprintf("%s%s?client=%s&secret=%s", oaUrl, oaTokenApi, oaClient, oaSecret)
		log.Debugf("Debug oa tokens api url:%s", queryUrl)

		resp, err := http.Get(queryUrl)

		if err != nil {
			log.Fatalf("query oa token api error!%s", err)
			return "", err
		}
		res, err := bodyResolve(resp)
		if err != nil {
				return "", err
		}
		log.Debugf("Debug oa token api returned:%d", res.Code)
		code := res.Code
		if code != consts.OAOkCode {
			log.Fatalf("Getting token error!%s", res.Message)
			return "", errors.New(res.Message)
		}
		//data := res.Data
		//logging.Printf("Debug oa tokens api return:%d, %s", code, data.Token)
		token = res.Data.Token
		log.Debugf("Debug oa token:%s", token)
		return token, nil

}
