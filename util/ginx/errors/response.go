package errors

import (
	"basic-frame/util/common"
	"basic-frame/util/i18n/Localizer"
	"fmt"
	"github.com/pkg/errors"
)

// UnWrapResponse 解包响应错误
func UnWrapResponse(err error) *common.ResponseError {
	if v, ok := err.(*common.ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse 包装响应错误
func WrapResponse(err error, code, statusCode int, CallerNameAndLine, msg string, args ...interface{}) error {
	newErr := err

	if CallerNameAndLine != "" {
		newErr = errors.WithMessage(err, fmt.Sprintf("%s: %s ==>", CallerNameAndLine, msg))
	}
	res := &common.ResponseError{
		Code:       code,
		Message:    Localizer.I18n.Translate(msg, args...),
		ERR:        newErr,
		StatusCode: statusCode,
	}
	return res
}

// Wrap400Response 包装错误码为400的响应错误
func Wrap400Response(err error, CallerNameAndLine, msg string, args ...interface{}) error {
	return WrapResponse(err, 400, 400, CallerNameAndLine, msg, args...)
}

// Wrap500Response 包装错误码为500的响应错误
func Wrap500Response(err error, CallerNameAndLine, msg string, args ...interface{}) error {
	return WrapResponse(err, 500, 500, CallerNameAndLine, msg, args...)
}

// NewResponse 创建响应错误
func NewResponse(code, statusCode int, msg string, args ...interface{}) error {
	res := &common.ResponseError{
		Code:       code,
		Message:    Localizer.I18n.Translate(msg, args...),
		StatusCode: statusCode,
	}
	return res
}

// New400Response 创建错误码为400的响应错误
func New400Response(msg string, args ...interface{}) error {
	return NewResponse(400, 400, msg, args...)
}

// New500Response 创建错误码为500的响应错误
func New500Response(msg string, args ...interface{}) error {
	return NewResponse(500, 500, msg, args...)
}
