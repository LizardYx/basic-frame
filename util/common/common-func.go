package common

import (
	"fmt"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetAppPath 获取程序所在路径
func GetAppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	appPath := path[:index]
	return appPath
}

// GetFileList 获取目录下所有文件(ext:后缀名)
func GetFileList(path string, ext string) []string {
	var ss []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Create Path:", path)
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println("Create Path Error:", err.Error())
			return ss
		}
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ext) {
			ss = append(ss, f.Name())
		}
	}
	return ss
}

// CheckStringContain 检查字符串数组是否包含某个元素
func CheckStringContain(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// CheckIntContain 检查int数组是否包含某个元素
func CheckIntContain(ss []int, s int) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// GetUUID 获取随机唯一标识
func GetUUID() string {
	var err error
	return uuid.Must(uuid.NewV4(), err).String()
}

// Copy 结构体映射
func Copy(s, ts interface{}) error {
	return copier.Copy(ts, s)
}

// GetDBConnString 获取Mysql连接字符串
func (a *Mysql) GetDBConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		a.User, a.Password, a.Host, a.Port, a.DBName)
}

func (a *WebServer) GetWebServerAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
