package api

import (
	"github.com/WeiXinao/daily_your_go/app/pkg/options"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares"
	"github.com/WeiXinao/daily_your_go/gmicro/server/restserver/middlewares/auth"
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func newJWTAuth(opts *options.JwtOptions) (middlewares.AuthStrategy, error) {
	gjwt, err := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            opts.Realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(opts.Key),
		Timeout:          opts.Timeout,
		MaxRefresh:       opts.MaxRefresh,
		LogoutResponse: func(ctx *gin.Context, code int){
			ctx.JSON(code, nil)	
		},
		IdentityHandler:  claimHandlerFunc,
		IdentityKey:      middlewares.KeyUserID,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
	})
	if err!= nil {
		return nil, err
	}
	return auth.NewJWTStrategy(*gjwt), nil
}

func claimHandlerFunc(ctx *gin.Context) interface{} {
	claims := ginjwt.ExtractClaims(ctx)
	ctx.Set(middlewares.KeyUserID, claims["ID"])
	return claims[ginjwt.IdentityKey]
}
