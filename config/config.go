package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var configFile *appConfig

func getExcPath() string {
	file, _ := exec.LookPath(os.Args[0])
	// 获取包含可执行文件名称的路径
	path, _ := filepath.Abs(file)
	// 获取可执行文件所在目录
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return strings.Replace(ret, "\\", "/", -1)
}

func GetConfig() appConfig {
	if configFile != nil {
		return *configFile
	}
	path := getExcPath()
	var mode modeConfig = readJson[modeConfig](path + "/config/app.json")

	appConfig := readJson[appConfig](path + "/config/app." + mode.Mode + ".json")
	configFile = &appConfig
	return *configFile
}

func readJson[T any](path string) T {
	file, err := os.Open(path)
	if err != nil {
		panic(path + " config not found")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)

	var jsonF T
	decoder.Decode(&jsonF)
	return jsonF
}

type modeConfig struct {
	Mode string
}

type appConfig struct {
	DBUser          dbConfig
	Redis           redisConfig
	CodePush        codePush
	UrlPrefix       string
	Port            string
	ResourceUrl     string
	TokenExpireTime int64
}
type dbConfig struct {
	Write           dbConfigObj
	MaxIdleConns    uint
	MaxOpenConns    uint
	ConnMaxLifetime uint
}
type dbConfigObj struct {
	UserName string
	Password string
	Host     string
	Port     uint
	DBname   string
}
type redisConfig struct {
	Host     string
	Port     uint
	DBIndex  uint
	UserName string
	Password string
}
type codePush struct {
	FileLocal string
	Local     localConfig
	Aws       awsConfig
	Ftp       ftpConfig
}
type awsConfig struct {
	Endpoint         string
	Region           string
	S3ForcePathStyle bool
	KeyId            string
	Secret           string
	Bucket           string
}
type ftpConfig struct {
	ServerUrl string
	UserName  string
	Password  string
}
type localConfig struct {
	SavePath string
}
