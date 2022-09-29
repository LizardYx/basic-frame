package logger

import (
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

var (
	Log = logrus.New()
)

// Init 日志模块初始化函数
func Init() {
	// 设置日志等级
	setLogLevel()
	// 设置日志输出格式
	Log.SetFormatter(&logrus.TextFormatter{})
	// 日志分割
	allWriter, err := rotatelogs.New(
		path.Join("log", "basic-frame"+".%Y%m%d"),
		rotatelogs.WithLinkName(path.Join("log", "basic-frame")), // 为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithRotationTime(time.Hour*24),                // 设置日志分割的时间
		rotatelogs.WithMaxAge(time.Hour*time.Duration(72)),       // 设置文件清理前的最长保存时间
	)
	if err != nil {
		fmt.Println("Init log error: ", err.Error())
		os.Exit(1)
	}
	writers := []io.Writer{allWriter, os.Stdout}
	Log.SetOutput(io.MultiWriter(writers...))
}

func setLogLevel() {
	switch common.SysConfig.RunMode {
	case consts.RunModeDev:
		Log.SetLevel(5)
	case consts.RunModeTest:
		Log.SetLevel(5)
	case consts.RunModePro:
		// 只记录 error/fatal/panic 错误
		Log.SetLevel(2)
	default:
		Log.SetLevel(5)
	}
}
