package utils

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var key = []byte("8x&*i}.r")

func CreateToken(str string) string {
	b, err := desEncrypt([]byte(str), key)
	if err != nil {
		log.Panic(err.Error())
	}

	return base64.StdEncoding.EncodeToString(b)
}
func GetDecToken(str string) string {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		log.Panic(err.Error())
	}
	b, err = desDecrypt(b, key)
	if err != nil {
		log.Panic(err.Error())
	}
	return string(b)
}
func desEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = pKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func desDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = pKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func pKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

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
