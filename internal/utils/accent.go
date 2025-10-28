package utils

import (
	"regexp"
	"strings"
)

func ConvertAccented(text string) string {
	text = strings.ToLower(text)
	var RegexpA = `à|á|ạ|ã|ả|ă|ắ|ằ|ẳ|ẵ|ặ|â|ấ|ầ|ẩ|ẫ|ậ`
	var RegexpE = `è|ẻ|ẽ|é|ẹ|ê|ề|ể|ễ|ế|ệ`
	var RegexpI = `ì|ỉ|ĩ|í|ị`
	var RegexpU = `ù|ủ|ũ|ú|ụ|ư|ừ|ử|ữ|ứ|ự`
	var RegexpY = `ỳ|ỷ|ỹ|ý|ỵ`
	var RegexpO = `ò|ỏ|õ|ó|ọ|ô|ồ|ổ|ỗ|ố|ộ|ơ|ờ|ở|ỡ|ớ|ợ`
	var RegexpD = `đ`
	rega := regexp.MustCompile(RegexpA)
	rege := regexp.MustCompile(RegexpE)
	regi := regexp.MustCompile(RegexpI)
	rego := regexp.MustCompile(RegexpO)
	regu := regexp.MustCompile(RegexpU)
	regy := regexp.MustCompile(RegexpY)
	regd := regexp.MustCompile(RegexpD)
	text = rega.ReplaceAllLiteralString(text, "a")
	text = rege.ReplaceAllLiteralString(text, "e")
	text = regi.ReplaceAllLiteralString(text, "i")
	text = rego.ReplaceAllLiteralString(text, "o")
	text = regu.ReplaceAllLiteralString(text, "u")
	text = regy.ReplaceAllLiteralString(text, "y")
	text = regd.ReplaceAllLiteralString(text, "d")

	return text
}
