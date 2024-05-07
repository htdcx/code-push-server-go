package model

import (
	"com.lc.go.codepush/server/db"
	"com.lc.go.codepush/server/model/constants"
	"gorm.io/gorm"
)

var userDb, _ = db.GetUserDB()

func queryPage[T any](db *gorm.DB, page constants.PageBean) *constants.PageData[T] {
	var datas []T
	var total int64
	db.Count(&total)
	err := db.Offset(page.Page * page.Rows).Limit(page.Rows).Find(&datas).Error
	if err != nil {
		return nil
	}
	return &constants.PageData[T]{
		TotalCount: total,
		Data:       datas,
	}
}

func Update[T any](saveData *T) {
	userDb.Updates(&saveData)
}

func Create[T any](saveData *T) (err error) {
	tx := userDb.Create(&saveData)
	return tx.Error
}

func GetOne[T any](sql string, arg any) *T {
	var t *T
	err := userDb.Select("*").Where(sql, arg).First(&t).Error
	if err != nil {
		return nil
	}
	return t
}

func GetList[T any](sql string, arg any) *[]T {
	var t *[]T
	err := userDb.Select("*").Where(sql, arg).Find(&t).Error
	if err != nil {
		return nil
	}
	return t
}

func Delete[T any](deleteData T) error {
	return userDb.Delete(deleteData).Error
}

func DeleteWhere(query string, args string, deleteData any) error {
	return userDb.Where(query, args).Delete(deleteData).Error
}
