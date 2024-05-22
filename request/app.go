package request

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"com.lc.go.codepush/server/config"
	"com.lc.go.codepush/server/db"
	"com.lc.go.codepush/server/db/redis"
	"com.lc.go.codepush/server/model"
	"com.lc.go.codepush/server/model/constants"
	"com.lc.go.codepush/server/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jlaffaye/ftp"
	"gorm.io/gorm"
)

type App struct{}

type createAppReq struct {
	AppName *string `json:"appName" binding:"required"`
	OS      *int    `json:"os" binding:"required"`
}

func (App) CreateApp(ctx *gin.Context) {
	createAppInfo := createAppReq{}
	if err := ctx.ShouldBindBodyWith(&createAppInfo, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)
		oldApp := model.App{}.GetAppByUidAndAppName(uid, *createAppInfo.AppName)
		if oldApp != nil {
			log.Panic("AppName " + *createAppInfo.AppName + " exist")
		}
		if *createAppInfo.OS != 1 && *createAppInfo.OS != 2 {
			log.Panic("OS error")
		}
		newApp := model.App{
			Uid:        &uid,
			AppName:    createAppInfo.AppName,
			OS:         createAppInfo.OS,
			CreateTime: utils.GetTimeNow(),
		}
		model.Create[model.App](&newApp)
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		log.Panic(err.Error())
	}
}

type createBundleReq struct {
	AppName     *string `json:"appName" binding:"required"`
	Deployment  *string `json:"deployment" binding:"required"`
	DownloadUrl *string `json:"downloadUrl" binding:"required"`
	Description *string `json:"description" binding:"required"`
	Version     *string `json:"version" binding:"required"`
	Size        *int64  `json:"size" binding:"required"`
	Hash        *string `json:"hash" binding:"required"`
}

func (App) CreateBundle(ctx *gin.Context) {
	createBundleReq := createBundleReq{}
	if err := ctx.ShouldBindBodyWith(&createBundleReq, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)

		app := model.App{}.GetAppByUidAndAppName(uid, *createBundleReq.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		deployment := model.Deployment{}.GetByAppidAndName(*app.Id, *createBundleReq.Deployment)
		if deployment == nil {
			log.Panic("Deployment " + *createBundleReq.Deployment + " not found")
		}
		deploymentVersion := model.DeploymentVersion{}.GetByKeyDeploymentIdAndVersion(*deployment.Id, *createBundleReq.Version)
		if deploymentVersion == nil {
			versionNum := utils.FormatVersionStr(*createBundleReq.Version)
			deploymentVersion = &model.DeploymentVersion{
				DeploymentId: deployment.Id,
				AppVersion:   createBundleReq.Version,
				VersionNum:   &versionNum,
				CreateTime:   utils.GetTimeNow(),
			}
			model.Create[model.DeploymentVersion](deploymentVersion)

			if deployment.VersionId != nil {
				deploymentVersionOld := model.GetOne[model.DeploymentVersion]("id=?", *deployment.VersionId)
				if utils.FormatVersionStr(*deploymentVersionOld.AppVersion) < utils.FormatVersionStr(*createBundleReq.Version) {
					deployment.VersionId = deploymentVersion.Id
					deployment.UpdateTime = utils.GetTimeNow()
					model.Update[model.Deployment](deployment)
				}
			} else {
				deployment.VersionId = deploymentVersion.Id
				deployment.UpdateTime = utils.GetTimeNow()
				model.Update[model.Deployment](deployment)
			}
		} else {
			nowPack := model.GetOne[model.Package]("id=?", deploymentVersion.CurrentPackage)
			if nowPack != nil && *nowPack.Hash == *createBundleReq.Hash {
				log.Panic("Upload package no modification")
			}
		}
		// uuid, _ := uuid.NewUUID()
		// hash := uuid.String()
		newPackage := model.Package{
			DeploymentId: deployment.Id,
			Size:         createBundleReq.Size,
			Hash:         createBundleReq.Hash,
			Download:     createBundleReq.DownloadUrl,
			Description:  createBundleReq.Description,
			Active:       utils.CreateInt(0),
			Installed:    utils.CreateInt(0),
			Failed:       utils.CreateInt(0),
			CreateTime:   utils.GetTimeNow(),
		}
		model.Create[model.Package](&newPackage)
		deploymentVersion.CurrentPackage = newPackage.Id
		deploymentVersion.UpdateTime = utils.GetTimeNow()
		model.Update[model.DeploymentVersion](deploymentVersion)
		redis.DelRedisObj(constants.REDIS_UPDATE_INFO + *deployment.Key + "*")
	} else {
		log.Panic(err.Error())
	}
}

type createDeploymentInfo struct {
	AppName        *string `json:"appName" binding:"required"`
	DeploymentName *string `json:"deploymentName" binding:"required"`
}

func (App) CreateDeployment(ctx *gin.Context) {
	createDeploymentInfo := createDeploymentInfo{}
	if err := ctx.ShouldBindBodyWith(&createDeploymentInfo, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)
		app := model.App{}.GetAppByUidAndAppName(uid, *createDeploymentInfo.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		deployment := model.Deployment{}.GetByAppidAndName(*app.Id, *createDeploymentInfo.DeploymentName)
		if deployment != nil {
			log.Panic("Deployment name " + *createDeploymentInfo.DeploymentName + " exist")
		}
		uuid, _ := uuid.NewUUID()
		key := uuid.String()
		newDeployment := model.Deployment{
			AppId:      app.Id,
			Name:       createDeploymentInfo.DeploymentName,
			Key:        &key,
			CreateTime: utils.GetTimeNow(),
		}
		err := model.Create[model.Deployment](&newDeployment)
		if err != nil {
			log.Panic(err.Error())
		}
		ctx.JSON(http.StatusOK, gin.H{
			"name": createDeploymentInfo.DeploymentName,
			"key":  key,
		})
	} else {
		log.Panic(err.Error())
	}
}
func (App) UploadBundle(ctx *gin.Context) {
	_, headers, err := ctx.Request.FormFile("file")
	if err != nil {
		log.Printf("Error when try to get file: %v", err)
	}
	file, err := headers.Open()
	if err != nil {
		log.Panic(err.Error())
	}
	defer file.Close()
	config := config.GetConfig()

	key := headers.Filename

	switch config.CodePush.FileLocal {
	case "local":
		exist := utils.Exists(config.CodePush.Local.SavePath)
		if !exist {
			err := os.MkdirAll(config.CodePush.Local.SavePath, 0777)
			if err != nil {
				log.Panic(err.Error())
			}
		}
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(file)
		if err != nil {
			log.Panic(err.Error())
		}
		os.WriteFile(config.CodePush.Local.SavePath+"/"+key, buf.Bytes(), 0777)
	case "aws":
		s3Config := &aws.Config{
			Credentials:      credentials.NewStaticCredentials(config.CodePush.Aws.KeyId, config.CodePush.Aws.Secret, ""),
			Endpoint:         aws.String(config.CodePush.Aws.Endpoint),
			Region:           aws.String(config.CodePush.Aws.Region),
			S3ForcePathStyle: aws.Bool(config.CodePush.Aws.S3ForcePathStyle),
		}
		newSession, _ := session.NewSession(s3Config)

		s3Client := s3.New(newSession)

		_, err = s3Client.PutObject(&s3.PutObjectInput{
			Body:   file,
			Bucket: aws.String(config.CodePush.Aws.Bucket),
			Key:    &key,
		})
		if err != nil {
			log.Panic(err.Error())
		}
	case "ftp":
		f, err := ftp.Dial(config.CodePush.Ftp.ServerUrl)
		if err != nil {
			log.Panic(err.Error())
		}
		err = f.Login(config.CodePush.Ftp.UserName, config.CodePush.Ftp.Password)
		if err != nil {
			log.Panic(err.Error())
		}

		err = f.Stor(key, file)
		if err != nil {
			log.Panic(err.Error())
		}
		if err := f.Quit(); err != nil {
			log.Panic(err.Error())
		}

	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

type lsDeploymentReq struct {
	ShowKey *bool   `json:"k" binding:"required"`
	AppName *string `json:"appName" binding:"required"`
}

type lsDeploymentInfo struct {
	AppName     *string           `json:"appName"`
	Deployments *[]deploymentInfo `json:"deployments"`
}

type deploymentInfo struct {
	DeploymentName *string `json:"deploymentName"`
	AppVersion     *string `json:"appVersion"`
	Active         *int    `json:"active"`
	Failed         *int    `json:"failed"`
	Installed      *int    `json:"installed"`
	DeploymentKey  *string `json:"deploymentKey"`
}

func (App) LsDeployment(ctx *gin.Context) {
	lsAppReq := lsDeploymentReq{}
	if err := ctx.ShouldBindBodyWith(&lsAppReq, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)
		app := model.App{}.GetAppByUidAndAppName(uid, *lsAppReq.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		var deploymentInfos []deploymentInfo
		deployment := model.Deployment{}.GetByAppids(*app.Id)

		for _, v := range *deployment {
			var key *string
			if *lsAppReq.ShowKey {
				key = v.Key
			}
			deploymentInfo := deploymentInfo{
				DeploymentName: v.Name,
				DeploymentKey:  key,
			}
			if v.VersionId != nil {
				deploymentVersion := model.GetOne[model.DeploymentVersion]("id=?", v.VersionId)
				deploymentInfo.AppVersion = deploymentVersion.AppVersion
				if deploymentVersion.CurrentPackage != nil {
					pack := model.GetOne[model.Package]("id=?", deploymentVersion.CurrentPackage)
					deploymentInfo.Active = pack.Active
					deploymentInfo.Failed = pack.Failed
					deploymentInfo.Installed = pack.Installed
				}
			}

			deploymentInfos = append(deploymentInfos, deploymentInfo)
		}
		lsAppInfo := lsDeploymentInfo{
			AppName:     app.AppName,
			Deployments: &deploymentInfos,
		}
		ctx.JSON(http.StatusOK, lsAppInfo)
	} else {
		log.Panic(err.Error())
	}
}

func (App) LsApp(ctx *gin.Context) {
	uid := ctx.MustGet(constants.GIN_USER_ID).(int)
	apps := model.GetList[model.App]("uid=?", uid)
	if len(*apps) <= 0 {
		log.Panic("No app")
	}
	var appsRep []string

	for _, v := range *apps {
		appsRep = append(appsRep, *v.AppName)
	}
	ctx.JSON(http.StatusOK, appsRep)
}

type checkBundleReq struct {
	AppName    *string `json:"appName" binding:"required"`
	Deployment *string `json:"deployment" binding:"required"`
	Version    *string `json:"version" binding:"required"`
}

func (App) CheckBundle(ctx *gin.Context) {
	checkBundleReq := checkBundleReq{}
	if err := ctx.ShouldBindBodyWith(&checkBundleReq, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)

		app := model.App{}.GetAppByUidAndAppName(uid, *checkBundleReq.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		deployment := model.Deployment{}.GetByAppidAndName(*app.Id, *checkBundleReq.Deployment)
		if deployment == nil {
			log.Panic("Deployment " + *checkBundleReq.Deployment + " not found")
		}
		var hash *string
		if deployment.VersionId != nil {
			deployment := model.DeploymentVersion{}.GetByKeyDeploymentIdAndVersion(*deployment.Id, *checkBundleReq.Version)
			if deployment != nil {
				pack := model.GetOne[model.Package]("id", deployment.CurrentPackage)
				hash = pack.Hash
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"appName": app.AppName,
			"os":      app.OS,
			"hash":    hash,
		})
	} else {
		log.Panic(err.Error())
	}
}

type delAppInfo struct {
	AppName *string `json:"appName" binding:"required"`
}

func (App) DelApp(ctx *gin.Context) {
	delAppInfo := delAppInfo{}
	if err := ctx.ShouldBindBodyWith(&delAppInfo, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)

		app := model.App{}.GetAppByUidAndAppName(uid, *delAppInfo.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		deployment := model.Deployment{}.GetByAppids(*app.Id)
		if deployment != nil && len(*deployment) > 0 {
			log.Panic("App exist deployment,Delete the deployment first and then delete the app ")
		}
		model.Delete[model.App](model.App{Id: app.Id})
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		log.Panic(err.Error())
	}
}

type delDeploymentInfo struct {
	AppName    *string `json:"appName" binding:"required"`
	Deployment *string `json:"deployment" binding:"required"`
}

func (App) DelDeployment(ctx *gin.Context) {
	delDeploymentInfo := delDeploymentInfo{}
	if err := ctx.ShouldBindBodyWith(&delDeploymentInfo, binding.JSON); err == nil {
		uid := ctx.MustGet(constants.GIN_USER_ID).(int)

		app := model.App{}.GetAppByUidAndAppName(uid, *delDeploymentInfo.AppName)
		if app == nil {
			log.Panic("App not found")
		}
		deployment := model.Deployment{}.GetByAppidAndName(*app.Id, *delDeploymentInfo.Deployment)
		if deployment == nil {
			log.Panic("Deployment " + *delDeploymentInfo.Deployment + " not found")
		}
		userDb, _ := db.GetUserDB()
		err := userDb.Transaction(func(tx *gorm.DB) error {
			if err := tx.Delete(model.Deployment{Id: deployment.Id}).Error; err != nil {
				panic("DeleteError:" + err.Error())
			}
			if err := tx.Where("deployment_id", *deployment.Id).Delete(model.DeploymentVersion{}).Error; err != nil {
				panic("DeleteError:" + err.Error())
			}
			if err := tx.Where("deployment_id", *deployment.Id).Delete(model.Package{}).Error; err != nil {
				panic("DeleteError:" + err.Error())
			}
			return nil
		})
		if err != nil {
			panic("DeleteError:" + err.Error())
		}

		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	} else {
		log.Panic(err.Error())
	}
}
