package auth

import (
	"errors"
	"fmt"
	"github.com/sunnywalden/sync-data/pkg/logging"
	"io/ioutil"
	"net/http"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/consts"
	"github.com/sunnywalden/sync-data/pkg/jsoner"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	log = logging.GetLogger()
)

func bodyResolve(resp *http.Response) (types.OAToken, error) {
	status := resp.StatusCode
	if status == 200 {
		defer resp.Body.Close()
		bodyC, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Debug api response:%s\n", string(bodyC))
		jsonMap, err := jsoner.JsonDecode(bodyC)
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
		log.Printf("Debug oa token api url:%s\n", queryUrl)

		resp, err := http.Get(queryUrl)

		if err != nil {
			log.Fatalf("query oa token api error!%s\n", err)
			return "", err
		}
		res, err := bodyResolve(resp)
		if err != nil {
				return "", err
		}
		log.Printf("Debug oa token api return:%d\n", res.Code)
		code := res.Code
		if code != consts.OAOkCode {
			log.Fatalf("Getting token error!%s\n", res.Message)
			return "", errors.New(res.Message)
		}
		data := res.Data
		//logging.Printf("Debug oa token api return:%d, %s\n", code, data.Token)
		token = data.Token
		return token, nil

}