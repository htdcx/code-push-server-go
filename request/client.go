package request

import (
	"log"
	"net/http"
	"strconv"

	"com.lc.go.codepush/server/config"
	"com.lc.go.codepush/server/db/redis"
	"com.lc.go.codepush/server/model"
	"com.lc.go.codepush/server/model/constants"
	"com.lc.go.codepush/server/utils"
	"github.com/gin-gonic/gin"
)

type Client struct{}
type updateInfo struct {
	DownloadUrl string `json:"download_url"`
	// Description            string `json:"description"`
	IsAvailable            bool   `json:"is_available"`
	IsDisabled             bool   `json:"is_disabled"`
	TargetBinaryRange      string `json:"target_binary_range"`
	PackageHash            string `json:"package_hash"`
	Label                  string `json:"label"`
	PackageSize            int64  `json:"package_size"`
	UpdateAppVersion       bool   `json:"update_app_version"`
	ShouldRunBinaryVersion bool   `json:"should_run_binary_version"`
	IsMandatory            bool   `json:"is_mandatory"`
}
type updateInfoRedisInfo struct {
	updateInfo
	NewVersion string
}

func (Client) CheckUpdate(ctx *gin.Context) {
	deploymentKey := ctx.Query("deployment_key")
	appVersion := ctx.Query("app_version")
	packageHash := ctx.Query("package_hash")
	// label := ctx.Query("label")
	// clientUniqueId := ctx.Query("client_unique_id")
	redisKey := constants.REDIS_UPDATE_INFO + deploymentKey + ":" + appVersion
	updateInfoRedis := redis.GetRedisObj[updateInfoRedisInfo](redisKey)
	updateInfo := updateInfo{}
	config := config.GetConfig()

	if updateInfoRedis == nil {
		updateInfoRedis = &updateInfoRedisInfo{}
		deployment := model.GetOne[model.Deployment]("key", deploymentKey)
		if deployment == nil {
			log.Panic("Key error")
		}
		deploymentVersion := model.DeploymentVersion{}.GetByKeyDeploymentIdAndVersion(*deployment.Id, appVersion)
		if deploymentVersion != nil {
			packag := model.GetOne[model.Package]("id", deploymentVersion.CurrentPackage)
			if packag != nil {
				// && *packag.Hash != packageHash
				updateInfoRedis.TargetBinaryRange = *deploymentVersion.AppVersion
				updateInfoRedis.PackageHash = *packag.Hash
				updateInfoRedis.PackageSize = *packag.Size
				updateInfoRedis.IsAvailable = true
				updateInfoRedis.IsMandatory = true
				label := strconv.Itoa(*packag.Id)
				updateInfoRedis.Label = label
				updateInfoRedis.DownloadUrl = config.ResourceUrl + *packag.Download
			}
		}
		deploymentVersionNew := model.DeploymentVersion{}.GetNewVersionByKeyDeploymentId(*deployment.Id)
		if deploymentVersionNew != nil {
			updateInfoRedis.NewVersion = *deploymentVersionNew.AppVersion
		}
		redis.SetRedisObj(redisKey, updateInfoRedis, -1)
	}
	if updateInfoRedis.PackageHash != "" {
		if updateInfoRedis.PackageHash != packageHash && appVersion == updateInfoRedis.TargetBinaryRange {
			updateInfo.TargetBinaryRange = updateInfoRedis.TargetBinaryRange
			updateInfo.PackageHash = updateInfoRedis.PackageHash
			updateInfo.PackageSize = updateInfoRedis.PackageSize
			updateInfo.IsAvailable = true
			updateInfo.IsMandatory = true
			updateInfo.Label = updateInfoRedis.Label
			updateInfo.DownloadUrl = updateInfoRedis.DownloadUrl
		} else if updateInfoRedis.NewVersion != "" && appVersion != updateInfoRedis.NewVersion && utils.FormatVersionStr(appVersion) < utils.FormatVersionStr(updateInfoRedis.NewVersion) {
			updateInfo.TargetBinaryRange = updateInfoRedis.NewVersion
			updateInfo.UpdateAppVersion = true
		}

	} else if updateInfoRedis.NewVersion != "" && appVersion != updateInfoRedis.NewVersion && utils.FormatVersionStr(appVersion) < utils.FormatVersionStr(updateInfoRedis.NewVersion) {
		updateInfo.TargetBinaryRange = updateInfoRedis.NewVersion
		updateInfo.UpdateAppVersion = true
	}

	ctx.JSON(http.StatusOK, gin.H{
		"update_info": updateInfo,
	})

}

type reportStatuReq struct {
	AppVersion                *string `json:"app_version"`
	DeploymentKey             *string `json:"deployment_key"`
	ClientUniqueId            *string `json:"client_unique_id"`
	Label                     *string `json:"label"`
	Status                    *string `json:"status"`
	PreviousLabelOrAppVersion *string `json:"previous_label_or_app_version"`
	PreviousDeploymentKey     *string `json:"previous_deployment_key"`
}

func (Client) ReportStatus(ctx *gin.Context) {
	json := reportStatuReq{}
	ctx.BindJSON(&json)
	if json.Status != nil {
		pack := model.GetOne[model.Package]("id=?", json.Label)
		if pack != nil {
			if *json.Status == "DeploymentSucceeded" {
				model.Package{}.AddActive(*pack.Id)
			} else if *json.Status == "DeploymentFailed" {
				model.Package{}.AddFailed(*pack.Id)
			}
		}

	}

	ctx.String(http.StatusOK, "OK")
}

type downloadReq struct {
	ClientUniqueId *string `json:"client_unique_id"`
	DeploymentKey  *string `json:"deployment_key"`
	Label          *string `json:"label"`
}

func (Client) Download(ctx *gin.Context) {
	json := downloadReq{}
	ctx.BindJSON(&json)
	pack := model.GetOne[model.Package]("id=?", json.Label)
	if pack != nil {
		model.Package{}.AddInstalled(*pack.Id)
	}
	ctx.String(http.StatusOK, "OK")
}
