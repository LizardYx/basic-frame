package errors

import (
	"basic-frame/util/i18n/Localizer"
	"github.com/pkg/errors"
)

// 定义别名
var (
	Wrap         = errors.Wrap // 同时附加堆栈和信息
	Wrapf        = errors.Wrapf
	WithStack    = errors.WithStack // 只附加调用堆栈信息
	WithMessagef = errors.WithMessagef
)

// New 新生成一个错误, 带堆栈信息
func New(errMessage string, args ...interface{}) error {
	if errMessage != "" {
		errMessage = Localizer.I18n.Translate(errMessage, args...)
	}
	return errors.New(errMessage)
}

// WithMessage 只附加新的错误信息
func WithMessage(err error, errMessage string, args ...interface{}) error {
	if errMessage != "" {
		errMessage = Localizer.I18n.Translate(errMessage, args...)
	}
	return errors.WithMessage(err, errMessage)
}
