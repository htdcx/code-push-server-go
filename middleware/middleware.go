package middleware

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"com.lc.go.codepush/server/model"
	"com.lc.go.codepush/server/model/constants"
	"com.lc.go.codepush/server/utils"
	"github.com/gin-gonic/gin"
)

// 檢查token
func CheckToken(ctx *gin.Context) {
	var token, _ = ctx.Cookie("token")
	if token == "" {
		token = ctx.GetHeader("token")
	}

	if token == "" {
		log.Panic("Token不能为空")
	}
	str := utils.GetDecToken(token)
	info := strings.Split(str, ":")
	expireTime, _ := strconv.ParseInt(info[1], 10, 64)
	if *utils.GetTimeNow() > expireTime {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 1100,
			"msg":  "Token expire",
		})
		ctx.Abort()
	} else {
		tokenNow := model.GetOne[model.Token]("token=?", token)
		if (tokenNow != nil && tokenNow.Del != nil && *tokenNow.Del) || tokenNow == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code": 1100,
				"msg":  "Token expire",
			})
			ctx.Abort()
		}
	}

	ctx.Set(constants.GIN_USER_ID, info[0])
}

// 異常處理
func Recover(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "*")
	lang := c.GetHeader("Accept-Language")
	c.Set(constants.GIN_LANG, lang)
	// 加载defer异常处理
	defer func() {
		if err := recover(); err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("Error:%s", err)
			// 返回统一的Json风格
			var msgStr string
			if fmt.Sprint(reflect.TypeOf(err)) == "string" {
				msgStr = fmt.Sprint(err)
			} else {
				msgStr = "system error"
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"msg":     msgStr,
				"success": false,
			})
			//终止后续操作
			c.Abort()
		}
	}()
	if c.Request.Method == "OPTIONS" {
		c.Writer.WriteHeader(http.StatusNoContent)
		c.Abort()
		// return
	}
	//继续操作
	c.Next()
}
