package model

type User struct {
	Id       *int    `gorm:"primarykey;autoIncrement;size:32"`
	UserName *string `gorm:"size:200" json:"userName"`
	Password *string `gorm:"size:45" json:"-"`
}
