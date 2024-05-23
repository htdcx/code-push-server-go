package main

import (
	"fmt"

	"com.lc.go.codepush/server/config"
	"com.lc.go.codepush/server/middleware"
	"com.lc.go.codepush/server/request"

	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("code-push-server-go V1.0.2")
	// gin.SetMode(gin.ReleaseMode)
	g := gin.Default()
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(middleware.Recover)
	configs := config.GetConfig()

	g.GET("/v0.1/public/codepush/update_check", request.Client{}.CheckUpdate)
	g.POST("/v0.1/public/codepush/report_status/deploy", request.Client{}.ReportStatus)
	g.POST("/v0.1/public/codepush/report_status/download", request.Client{}.Download)

	apiGroup := g.Group(configs.UrlPrefix)
	{
		apiGroup.POST("/login", request.User{}.Login)
	}
	authApi := apiGroup.Use(middleware.CheckToken)
	{
		authApi.POST("/createApp", request.App{}.CreateApp)
		authApi.POST("/createDeployment", request.App{}.CreateDeployment)
		authApi.POST("/createBundle", request.App{}.CreateBundle)
		authApi.POST("/checkBundle", request.App{}.CheckBundle)
		authApi.POST("/delApp", request.App{}.DelApp)
		authApi.POST("/delDeployment", request.App{}.DelDeployment)
		authApi.POST("/lsDeployment", request.App{}.LsDeployment)
		authApi.GET("/lsApp", request.App{}.LsApp)
		authApi.POST("/uploadBundle", request.App{}.UploadBundle)
	}

	g.Run(configs.Port)
}
