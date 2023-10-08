# basic-frame
包含国际化、日志、DB初始化、用户权限(RBAC3.0)等基本功能的项目基础框架代码

## i18n依赖
```shell
go install golang.org/x/text/cmd/gotext@latest
```

## swagger依赖
```shell
go install github.com/swaggo/swag/cmd/swag@latest
go get -u -v github.com/swaggo/gin-swagger
go get -u -v github.com/swaggo/files
```

## 生成/更新 swagger
```shell
make swagger
```

## 启动程序
```shell
make start
```

## 生成linux可执行文件
```shell
make linux-build
```

## 生成windows可执行文件
```shell
make windows-build
```

## 生成linux可执行程序(windows操作系统使用)
```shell
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build basic-frame
```