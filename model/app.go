package model

type App struct {
	Id         *int    `gorm:"primarykey;autoIncrement;size:32"`
	Uid        *int    `json:"uid"`
	AppName    *string `json:"appName"`
	OS         *int    `json:"os"`
	CreateTime *int64  `json:"createTime"`
}

func (App) GetAppByUidAndAppName(uid int, appName string) *App {
	var app *App
	err := userDb.Where("uid", uid).Where("app_name", appName).First(&app).Error
	if err != nil {
		return nil
	}
	return app
}
