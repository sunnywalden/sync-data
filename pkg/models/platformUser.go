package models

// platform user
type PlatUser struct {
	Id int `json:"id",gorm:"primary_key;id:int;not null;index:platuser_idx"`
	UserName string `bson:"username" json:"username",gorm:"username:varchar(256);not null;"`
	AuthKey string `bson:"authkey" json:"authkey",gorm:"authkey:varchar(256);not null;"`
}

type JwtToken struct {
	Token string `json:"token"`
}

