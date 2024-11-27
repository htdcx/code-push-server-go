package model

type User struct {
	Id       *int    `gorm:"primarykey;autoIncrement;size:32"`
	UserName *string `gorm:"size:200" json:"userName"`
	Password *string `gorm:"size:45" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (User) ChangePassword(uid int, password string) error {
	return userDb.Raw("update users set password=? where id=?", password, uid).Scan(&User{}).Error
}
