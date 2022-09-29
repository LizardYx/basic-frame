package translations

/**
@srclang:指定基础语言(BCP 47标签)
@update:需要执行的gotext中的函数
@out:消息输出目录所在路径(所在目录的相对路径)
@lang:需要翻译的语言,使用逗号分离
*/

//go:generate gotext -srclang=en_US update -out=catalog.go -lang=en_US,zh_CN basic-frame
