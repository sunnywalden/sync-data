package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sunnywalden/sync-data/pkg/errs"
	"github.com/sunnywalden/sync-data/pkg/types"
)

//TokenMiddleware, token check
func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var res = types.Response{
				Code: http.StatusUnauthorized,
				Msg: errs.ErrNotAuthorized.Error(),
				Data: nil,
		}

		_, err := VerifyToken(c)

		if err != nil {
			status := http.StatusUnauthorized
			res.Msg = "token invalid"
			c.JSON(
				status,
				res,
				)
			c.Abort()
		}
		c.Next()
	}
}
