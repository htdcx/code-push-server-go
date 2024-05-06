package model

type Token struct {
	Id         *int    `gorm:"primarykey;autoIncrement;size:32"`
	Uid        *int    `json:"uid"`
	Token      *string `json:"token"`
	ExpireTime *int64  `json:"expireTime"`
	Del        *bool   `json:"del"`
}

func (Token) TableName() string {
	return "token"
}
