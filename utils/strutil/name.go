package strutil

import (
	"unicode"
)

func GetNameStringLength(nameStr string) int {
	var strLength int
	for _, r := range nameStr {
		if unicode.Is(unicode.Han, r) { // 检查是否为中文字符
			strLength += 2
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) { // 检查是否为英文字符或数字
			strLength += 1
		}
	}
	return strLength
}
