package sensitiveword

import "testing"

func TestValidationPolicyAndInvalidUTF8(t *testing.T) {
	if IsStringValid("正常名字!") || IsStringValid("正常名字\xff") {
		t.Fatal("invalid name accepted")
	}
	if !IsStringValid("普通玩家123") {
		t.Fatal("ordinary name rejected")
	}
}
