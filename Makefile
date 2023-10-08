.PHONY: all

APP = BasicFrame
APP_SERVER_BIN = ./bin/${APP}

swagger:
	swag fmt -d ./ --exclude ./config,./docs,./middleware,./util & \
    swag init

start:
	go run basic-frame

linux-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${APP_SERVER_BIN} basic-frame

windows-build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${APP_SERVER_BIN}.exe basic-frame