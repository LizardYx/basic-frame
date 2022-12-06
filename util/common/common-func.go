package common

import (
	"basic-frame/util/consts"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"time"
	"unsafe"
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

// GetRandomString 获取随机字符串
func GetRandomString(l int) string {
	str := "0123456789AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	bytes := []byte(str)
	var result []byte = make([]byte, 0, l)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return BytesToString(result)
}

// BytesToString 0 拷贝转换 slice byte 为 string
func BytesToString(b []byte) (s string) {
	_bptr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	_sptr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	_sptr.Data = _bptr.Data
	_sptr.Len = _bptr.Len
	return s
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

// GetWebServerAddr 获取web服务地址
func (a *WebServer) GetWebServerAddr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

func (r *ResponseError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", r.ERR)
		}
		// fallthrough
	case 's':
		_, _ = io.WriteString(s, r.ERR.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", r.ERR.Error())
	}
}

func (a *SystemConfig) GetGinRunModel() string {
	runModel := gin.DebugMode
	if a.RunMode == consts.RunModeTest {
		runModel = gin.TestMode
	} else if a.RunMode == consts.RunModeRelease {
		runModel = gin.ReleaseMode
	}
	return runModel
}
