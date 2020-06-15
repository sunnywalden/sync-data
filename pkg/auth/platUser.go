package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/sunnywalden/sync-data/config"
	"github.com/sunnywalden/sync-data/pkg/databases"
	"github.com/sunnywalden/sync-data/pkg/errs"
	"github.com/sunnywalden/sync-data/pkg/models"
	"github.com/sunnywalden/sync-data/pkg/types"
)

var (
	ExpireTime = 3600 * 24     //token expire time
)

// GenerateToken, new token
func GenerateToken(user *models.PlatUser) (string, error) {
	claims := &types.JWTClaims{
		User: *user,
	}
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(ExpireTime)).Unix()

	authKey:= user.AuthKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(authKey))
}

func VerifyToken(c *gin.Context) (claims *types.JWTClaims, err error) {
	tokenStr := c.Request.Header.Get("X-Authorization-Token")
	if tokenStr == "" {
		err = errs.ErrTokenNotExists
		return nil, err
	}

	platUserName := c.PostForm("platuser")

	db, err := databases.Conn(config.Conf)
	if err != nil {
		return nil, err
	}

	var platUser models.PlatUser
	db = db.Where(&models.PlatUser{UserName: platUserName}).First(&platUser)

	token, err := jwt.ParseWithClaims(tokenStr, &types.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(platUser.AuthKey), nil
	})

	if claims, ok := token.Claims.(*types.JWTClaims); ok && token.Valid {
		fmt.Printf("%v %v", claims.Id, claims.StandardClaims.ExpiresAt)
		return claims, nil
	} else {
		fmt.Println(err)
		return nil, err
	}
}

// 更新token
func RefreshToken(c *gin.Context) (string, error) {
	claims, err := VerifyToken(c)
	if err != nil {
		db, err := databases.Conn(config.Conf)
		if err != nil {
			return "", err
		}

		platUserName := c.PostForm("platuser")

		var platUser models.PlatUser
		db = db.Where(&models.PlatUser{UserName: platUserName}).First(&platUser)
		return GenerateToken(&platUser)
	} else {
		return claims.Id, nil
	}
}
