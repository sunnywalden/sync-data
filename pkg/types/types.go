package types

import "github.com/sunnywalden/sync-data/pkg/models"

type OaToken string

type OAToken struct {
	Code    int     `json:"code"`
	Data    *OAAuth `json:"data"`
	Message string  `json:"msg"`
}

type OAAuth struct {
	ExpiredIn int64  `json:"expired_in(s)"`
	Token     OaToken `json:"token"`
}


type UserInfo struct {
	Code    int     `json:"code"`
	Data    []models.User `json:"data"`
	Message string  `json:"msg"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
