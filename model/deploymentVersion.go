package model

type DeploymentVersion struct {
	Id             *int    `gorm:"primarykey;autoIncrement;size:32"`
	DeploymentId   *int    `json:"deploymentId"`
	AppVersion     *string `json:"appVersion"`
	VersionNum     *int64  `json:"version_num"`
	CurrentPackage *int    `json:"currentPackage"`
	UpdateTime     *int64  `json:"updateTime"`
	CreateTime     *int64  `json:"createTime"`
}

func (DeploymentVersion) TableName() string {
	return "deployment_version"
}

func (DeploymentVersion) GetByKeyDeploymentIdAndVersion(deploymentId int, version string) *DeploymentVersion {
	var deploymentVersion *DeploymentVersion
	err := userDb.Where("deployment_id", deploymentId).Where("app_version", version).First(&deploymentVersion).Error
	if err != nil {
		return nil
	}
	return deploymentVersion
}

func (DeploymentVersion) GetNewVersionByKeyDeploymentId(deploymentId int) *DeploymentVersion {
	var deploymentVersion *DeploymentVersion
	err := userDb.Where("deployment_id", deploymentId).Order("version_num desc").First(&deploymentVersion).Error
	if err != nil {
		return nil
	}
	return deploymentVersion
}

func (DeploymentVersion) UpdateCurrentPackage(id int, pid *int) {
	userDb.Raw("update deployment_version set current_package=? where id=?", pid, id).Scan(&DeploymentVersion{})
}
