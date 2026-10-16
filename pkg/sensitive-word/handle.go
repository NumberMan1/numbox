package sensitiveword

import (
	_ "embed"
	"unicode"
	"unicode/utf8"

	"github.com/NumberMan1/numbox/pkg/matcher"
	"github.com/NumberMan1/numbox/pkg/utils"
)

var sensitiveMatcher *matcher.Matcher

//go:embed sensitive_word.json
var sensitiveWordsJson string

func init() {
	var sensitiveWords []string
	utils.UnmarshalFromJSON(sensitiveWordsJson, &sensitiveWords)
	sensitiveMatcher = matcher.NewMatcher(sensitiveWords)
}

func IsStringValid(str string) bool {
	if !utf8.ValidString(str) || sensitiveMatcher.Contains(str) {
		return false
	}
	for _, char := range str {
		if unicode.IsPunct(char) {
			return false
		}
	}
	return true
}
