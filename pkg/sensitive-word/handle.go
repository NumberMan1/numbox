package sensitiveword

import (
	_ "embed"
	"strings"
	"unicode"

	"github.com/NumberMan1/numbox/pkg/collection"
	"github.com/NumberMan1/numbox/pkg/utils"
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
