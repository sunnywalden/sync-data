package types

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
	Data    []User `json:"data"`
	Message string  `json:"msg"`
}

// User, oa user
type User struct {
	Id int `json:"departmentName",gorm:"primary_key;id:int;not null;index:user_idx"`
	Department string `json:"departmentName",gorm:"depart:varchar(256);not null;"`
	Title string `json:"jobTitle",gorm:"title:varchar(256);null;"`
	Name string `json:"lastName",gorm:"name:varchar(256);not null;"`
	LoginId string `json:"loginId",gorm:"loginId:varchar(256);not null;index:user_loginidx"`
	NickName string `json:"nickName",gorm:"nickname:varchar(256);not null;index:user_nicknamex"`
	Status uint8 `json:"status",gorm:"status:int;not null;"`
}
