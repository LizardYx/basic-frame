package errors

import (
	"github.com/pkg/errors"
)

// 定义别名
var (
	New          = errors.New  // 新生成一个错误, 带堆栈信息
	Wrap         = errors.Wrap // 同时附加堆栈和信息
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack   // 只附加调用堆栈信息
	WithMessage  = errors.WithMessage // 只附加新的信息
	WithMessagef = errors.WithMessagef
)

//func PrintCallerNameAndLine() string {
//	pc, _, line, _ := runtime.Caller(2)
//	return fmt.Sprintf("%s()@%s", runtime.FuncForPC(pc).Name(), strconv.Itoa(line))
//}
//
//func New(message string) error {
//	return errors.New(fmt.Sprintf("%s: %s", PrintCallerNameAndLine(), message))
//}
//
//func Wrap(err error, message string) error {
//	return errors.Wrap(err, fmt.Sprintf("%s: %s ==>", PrintCallerNameAndLine(), message))
//}
//
//func WithMessage(err error, message string) error {
//	return errors.WithMessage(err, fmt.Sprintf("%s: %s ==>", PrintCallerNameAndLine(), message))
//}
