package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/consts"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log = logging.GetLogger()
)

func jsonDecode(bodyC []byte) (types.OAToken, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var jsonMap types.OAToken

	err := json.Unmarshal(bodyC, &jsonMap)
	if err != nil {
		log.Fatalf("json decode err!%s\n", err.Error())
		return types.OAToken{}, err
	} else {
		return jsonMap, nil
	}
}

func bodyResolve(resp *http.Response) (types.OAToken, error) {
	status := resp.StatusCode
	if status == 200 {
		defer resp.Body.Close()
		bodyC, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Debug api response:%s\n", string(bodyC))
		jsonMap, err := jsonDecode(bodyC)
		if err != nil {
			log.Fatalf("Read api response error!%s\n", err)
			return types.OAToken{}, err
		}
		return jsonMap, nil
	}

	return types.OAToken{}, errors.New(resp.Status)
}

func GetToken() (token types.OaToken, err error) {
	oaConf := config.Conf.OA
	oaUrl := oaConf.OaUrl
	oaTokenApi := oaConf.TokenApi
	oaClient := oaConf.OaClient
	oaSecret := oaConf.OaSecret

		queryUrl := fmt.Sprintf("%s%s?client=%s&secret=%s", oaUrl, oaTokenApi, oaClient, oaSecret)
		log.Printf("Debug oa tokens api url:%s\n", queryUrl)

		resp, err := http.Get(queryUrl)

		if err != nil {
			log.Fatalf("query oa tokens api error!%s\n", err)
			return "", err
		}
		res, err := bodyResolve(resp)
		if err != nil {
				return "", err
		}
		log.Printf("Debug oa tokens api return:%d\n", res.Code)
		code := res.Code
		if code != consts.OAOkCode {
			log.Fatalf("Getting tokens error!%s\n", res.Message)
			return "", errors.New(res.Message)
		}
		data := res.Data
		//logging.Printf("Debug oa tokens api return:%d, %s\n", code, data.Token)
		token = data.Token
		return token, nil

}