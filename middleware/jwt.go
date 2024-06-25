package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-grain/gin-layout/config"
	consts "github.com/go-grain/gin-layout/const"
	"github.com/go-grain/gin-layout/internal/repo/system/query"
	model "github.com/go-grain/gin-layout/model/system"
	"github.com/go-grain/gin-layout/pkg/encrypt"
	jwtx "github.com/go-grain/gin-layout/pkg/jwt"
	redisx "github.com/go-grain/gin-layout/pkg/redis"
	"github.com/go-grain/gin-layout/pkg/response"
	"net/http"
	"time"
)

func JwtAuth(conf config.Config, rdb redisx.IRedis) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reply := response.Response{}
		jwt := jwtx.Jwt{}
		tokenString := ctx.GetHeader("G-Token")
		tokenClaims, err := jwt.ParseToken(tokenString, conf.JWT.SecretKey)
		if err != nil {
			reply.WithCode(http.StatusUnauthorized).WithMessage(err.Error()).Fail(ctx)
			ctx.Abort()
			return
		}

		black, _ := rdb.GetInt(fmt.Sprintf("%s%s", consts.TokenBlack, encrypt.MD5(tokenString)))
		switch black {
		case 120:
			reply.WithCode(http.StatusUnauthorized).WithMessage("账号进入黑名单列表,无法在继续为您服务").Fail(ctx)
			ctx.Abort()
			return
		case 110:
			reply.WithCode(http.StatusUnauthorized).WithMessage("登录异常,请重新登录").Fail(ctx)
			ctx.Abort()
			return
		case 100:
			reply.WithCode(http.StatusUnauthorized).WithMessage("无效请求 账号已退出登录").Fail(ctx)
			ctx.Abort()
			return
		}
		//获取用户信息
		sysUser := &model.SysUser{}
		if err = rdb.GetObject(consts.UserInfo+tokenClaims.Uid, sysUser); err != nil {
			sysUser, err = query.Q.SysUser.Where(query.SysUser.UID.Eq(tokenClaims.Uid)).First()
			_ = rdb.SetObject(consts.UserInfo+tokenClaims.Uid, sysUser, 180)
		}

		//把用户相关信息都塞到ctx去,方便下游使用
		if err == nil {
			ctx.Set("username", sysUser.Username)
			ctx.Set("nickname", sysUser.Nickname)
			ctx.Set("email", sysUser.Email)
			ctx.Set("mobil", sysUser.Mobile)
		}
		expired := int64(tokenClaims.ExpiresAt.Time.Sub(time.Now()).Seconds())
		ctx.Set("expTokenAt", expired)
		ctx.Set("uid", tokenClaims.Uid)
		ctx.Set("role", tokenClaims.Role)
		ctx.Set("token", encrypt.MD5(tokenString))
		ctx.Next()
	}
}
