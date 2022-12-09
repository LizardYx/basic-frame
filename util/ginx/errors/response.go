package errors

import (
	"basic-frame/util/common"
	"basic-frame/util/i18n/Localizer"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

// UnWrapResponse 解包响应错误
func UnWrapResponse(err error) *common.ResponseError {
	if v, ok := err.(*common.ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse 包装响应错误
func WrapResponse(err error, code, statusCode int, msg string, args ...interface{}) error {
	res := &common.ResponseError{
		Code:       code,
		Message:    Localizer.I18n.Translate(msg, args...),
		ERR:        errors.WithMessage(err, fmt.Sprintf(msg, args...)),
		StatusCode: statusCode,
	}
	return res
}

// Wrap400Response 包装错误码为400的响应错误
func Wrap400Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, http.StatusBadRequest, http.StatusBadRequest, msg, args...)
}

// Wrap500Response 包装错误码为500的响应错误
func Wrap500Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, http.StatusInternalServerError, http.StatusInternalServerError, msg, args...)
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

// NewIANAResponse 创建在IANA注册的错误
func NewIANAResponse(httpCode int, msg string, args ...interface{}) error {
	if msg != "" {
		return NewResponse(httpCode, httpCode, msg, args...)
	} else {
		return NewResponse(httpCode, httpCode, http.StatusText(httpCode))
	}
}
