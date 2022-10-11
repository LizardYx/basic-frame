package logger

import (
	"basic-frame/util/common"
	"basic-frame/util/consts"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"path"
	"strings"
	"time"
)

var (
	Log = logrus.New()
)

// InitLogger 日志模块初始化函数
func InitLogger() {
	// 设置日志等级
	setLogLevel(Log)

	// 设置日志输出格式
	Log.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 日志分割
	logName := strings.ToLower(common.SysConfig.AppName)
	rotationTime := common.SysConfig.Logger.RotationTime
	maxAge := common.SysConfig.Logger.MaxAge
	allWriter, err := rotatelogs.New(
		path.Join("log", logName+".%Y%m%d"),
		rotatelogs.WithLinkName(path.Join("log", logName)),                 // 为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithRotationTime(time.Duration(rotationTime)*time.Hour), // 设置日志分割的时间
		rotatelogs.WithMaxAge(time.Duration(maxAge)*time.Hour),             // 设置文件清理前的最长保存时间
	)
	if err != nil {
		log.Fatal("Init log error: ", err.Error())
	}

	// 设置日志输出到日志文件
	writers := []io.Writer{allWriter}
	Log.SetOutput(io.MultiWriter(writers...))
}

func setLogLevel(logger *logrus.Logger) {
	switch common.SysConfig.RunMode {
	case consts.RunModeDev:
		logger.SetLevel(logrus.DebugLevel)
	case consts.RunModeTest:
		logger.SetLevel(logrus.DebugLevel)
	case consts.RunModePro:
		// 只记录 error/fatal/panic 错误
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
	}
}
