package model

type Package struct {
	Id           *int    `gorm:"primarykey;autoIncrement;size:32"`
	DeploymentId *int    `json:"deploymentId"`
	Size         *int64  `json:"size"`
	Hash         *string `json:"hash"`
	Download     *string `json:"download"`
	Active       *int    `json:"active"`
	Failed       *int    `json:"failed"`
	Installed    *int    `json:"installed"`
	CreateTime   *int64  `json:"create_time"`
	Description  *string `json:"description"`
}

func (Package) TableName() string {
	return "package"
}

func (Package) AddActive(pid int) {
	userDb.Raw("update package set active=active+1 where id=?", pid).Scan(&Package{})
}

func (Package) AddFailed(pid int) {
	userDb.Raw("update package set failed=failed+1 where id=?", pid).Scan(&Package{})
}

func (Package) AddInstalled(pid int) {
	userDb.Raw("update package set installed=installed+1 where id=?", pid).Scan(&Package{})
}
