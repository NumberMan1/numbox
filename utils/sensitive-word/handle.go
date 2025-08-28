package sensitiveword

import (
	_ "embed"
	"github.com/NumberMan1/numbox/utils"
	"github.com/NumberMan1/numbox/utils/collection"
	"strings"
	"unicode"
)

var sensitiveWordsSet collection.Set[string]

//go:embed sensitive_word.json
var sensitiveWordsJson string

func init() {
	var sensitiveWords []string
	utils.UnmarshalFromJSON(sensitiveWordsJson, &sensitiveWords)
	sensitiveWordsSet = collection.NewSet(sensitiveWords...)
}

func IsStringValid(str string) bool {
	if sensitiveWordsSet.Contains(str) {
		return false
	}
	for word := range sensitiveWordsSet {
		if word != "" && strings.Contains(str, word) {
			return false
		}
	}
	for _, char := range str {
		if unicode.IsPunct(char) {
			return false
		}
	}
	return true
}
