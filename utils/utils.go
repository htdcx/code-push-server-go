package utils

import (
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetTimeNow() *int64 {
	t := time.Now().UnixMilli()
	return &t
}

func CreateInt(num int) *int {
	return &num
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func FormatVersionStr(v string) int64 {
	vs := strings.Split(v, ".")
	if len(vs) <= 0 {
		log.Panic("Version str error")
	}
	var vNum int64
	ReverseArr(vs)
	for index, v := range vs {
		num, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Panic(err.Error())
		}
		for i := 0; i < index; i++ {
			num = num * 100
		}
		vNum += num
	}
	return vNum
}
func ReverseArr(s interface{}) {
	sort.SliceStable(s, func(i, j int) bool {
		return true
	})
}
