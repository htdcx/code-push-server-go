package model

type Deployment struct {
	Id         *int    `gorm:"primarykey;autoIncrement;size:32"`
	AppId      *int    `json:"appId"`
	Name       *string `json:"name"`
	Key        *string `json:"key"`
	VersionId  *int    `json:"versionId"`
	UpdateTime *int64  `json:"updateTime"`
	CreateTime *int64  `json:"createTime"`
}

func (Deployment) TableName() string {
	return "deployment"
}

func (Deployment) GetByAppidAndName(appId int, name string) *Deployment {
	var deployment *Deployment
	err := userDb.Where("app_id", appId).Where("name", name).First(&deployment).Error
	if err != nil {
		return nil
	}
	return deployment
}

func (Deployment) GetByAppids(appId int) *[]Deployment {
	var deployment *[]Deployment
	err := userDb.Where("app_id", appId).Find(&deployment).Error
	if err != nil {
		return nil
	}
	return deployment
}
