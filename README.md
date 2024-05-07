# code-push-server-go
Codepush server go is compatible with [react-native-code-push](https://github.com/microsoft/react-native-code-push). Need to be used with [code-push-go](https://github.com/htdcx/code-push-go). Only supported react-native

## Support version
- [mysql](https://dev.mysql.com/downloads/mysql/)  >= 8.0
- [golang](https://go.dev/dl/) >= 1.21.5
- [redis](https://redis.io/downloads/)  >= 5.0

## Support client version
- [react-native-code-push](https://github.com/microsoft/react-native-code-push) >= 7.0

## Support storage
- Local
- AWS S3 
- FTP

## Before installation, please ensure that the following procedures have been installed
- mysql
- golang
- redis

## Install code-push-server
```shell
git clone https://github.com/htdcx/code-push-server-go.git
cd code-push-server-go
import code-push.sql to mysql
```
### Configuration mysql,redis,storage
``` shell
cd config
vi (app.json or app.dev.json or app.prod.json) 
# app.json
{
    "mode":"prod" #run read config app.{mode}.json
}
# app.prod.json
{
    "DBUser": {
        "Write": {
            "UserName": "",
            "Password": "",
            "Host": "127.0.0.1",
            "Port": 3306,
            "DBname": ""
        },
        "MaxIdleConns": 10,
        "MaxOpenConns": 100,
        "ConnMaxLifetime": 1
    },
    "Redis": {
        "Host": "127.0.0.1",
        "Prot": 6379,
        "DBIndex": 0,
        "UserName": "",
        "Password": ""
    },
    "CodePush": {
        "FileLocal":(local,aws,ftp),
        "Local":{
            "SavePath":"./bundels"
        },
        "Aws":{
            "Endpoint":"",
            "Region":"",
            "S3ForcePathStyle":true,
            "KeyId":"",
            "Secret":"",
            "Bucket":""
        },
        "Ftp":{
            "ServerUrl":"",
            "UserName":"",
            "Password":""
        }
    },
    "UrlPrefix": "/",
    "ResourceUrl": (nginx config url or s3),
    "Port": ":8080",
    "TokenExpireTime": 30 (day)
}

```
#### Build
``` shell
#MacOS pack GOOS:windows,darwin
CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -o code-push-server-go(.exe) main.go

#Windows pack
set GOARCH=amd64
set GOOS=linux #windows,darwin
go build -o code-push-server-go(.exe) main.go

#copy config/app.(model).json and config/app.json to run dir 

#Linux server
chmod +x code-push-go

#run
./code-push-server-go
```
### Default user name and password
- Username:admin
- Password:admin

### Change password and user name
- Change mysql users tables (password need md5)

### Use [code-push-go](https://github.com/htdcx/code-push-go)
``` shell
./code-push-go login -u (userName) -p (password) -h (serverUrl)
```
### Configuration client [react-native-code-push](https://github.com/microsoft/react-native-code-push)

``` shell
#ios add to Info.plist
<key>CodePushServerURL</key>
<string>${CODE_PUSH_SERVER_URL}</string>

#android add to res/value/strings.xml
<string moduleConfig="true" name="CodePushServerUrl">${CODE_PUSH_SERVER_URL}</string>
```

## Developing
- [ ] Delete app
- [ ] Delete deployment
- [ ] Rollback bundel

## License
MIT License [Read](https://github.com/htdcx/code-push-server-go/blob/main/LICENSE)
