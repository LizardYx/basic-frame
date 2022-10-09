package Localizer

import (
	"basic-frame/util/consts"
	_ "basic-frame/util/i18n/translations"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	I18n = new(Localizer)
)

// localizerNames 可用的语言列表
var localizerNames = []string{consts.LanguageZh, consts.LanguageEn}

// Localizer 国际化对象
type Localizer struct {
	Name    string
	printer *message.Printer
}

// GetLocalizer 获取指定的国际化对象
func GetLocalizer(name string) *Localizer {
	var localizerName string

	for index, _ := range localizerNames {
		if localizerNames[index] == name {
			localizerName = name
			break
		}
		if index == (len(localizerNames) - 1) {
			// 如果未找到相应的语言，则使用系统默认语言
			localizerName = consts.DefaultLang
		}
	}
	return &Localizer{
		Name:    localizerName,
		printer: message.NewPrinter(language.MustParse(localizerName)),
	}
}

// Translate 文本翻译
func (a Localizer) Translate(key message.Reference, args ...interface{}) string {
	return a.printer.Sprintf(key, args...)
}
