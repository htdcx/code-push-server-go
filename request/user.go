package request

import (
	"log"
	"net/http"

	"com.lc.go.codepush/server/config"
	"com.lc.go.codepush/server/model"
	"com.lc.go.codepush/server/model/constants"
	"com.lc.go.codepush/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

type User struct{}

type loginUser struct {
	UserName *string `json:"userName" binding:"required"`
	Password *string `json:"password" binding:"required"`
}

func (User) Login(ctx *gin.Context) {
	loginUser := loginUser{}
	if err := ctx.ShouldBindBodyWith(&loginUser, binding.JSON); err == nil {
		user := model.GetOne[model.User]("user_name", &loginUser.UserName)
		if user == nil || *user.Password != *loginUser.Password {
			panic("UserName or Psssword error")
		}
		uuid, _ := uuid.NewUUID()
		timeNow := utils.GetTimeNow()
		expireTime := *timeNow + (config.GetConfig().TokenExpireTime * 24 * 60 * 60 * 1000)
		token := uuid.String()
		del := false
		tokenInfo := model.Token{
			Uid:        user.Id,
			Token:      &token,
			ExpireTime: &expireTime,
			Del:        &del,
		}
		err := model.Create[model.Token](&tokenInfo)
		if err != nil {
			panic("create token error")
		}
		ctx.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	} else {
		log.Panic(err.Error())
	}
}

type changePasswordReq struct {
	Password *string `json:"password" binding:"required"`
}

func (User) ChangePassword(ctx *gin.Context) {
	req := changePasswordReq{}
	if err := ctx.ShouldBindBodyWith(&req, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)
		err := model.User{}.ChangePassword(uid, *req.Password)
		if err != nil {
			panic(err.Error())
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		log.Panic(err.Error())
	}
}
