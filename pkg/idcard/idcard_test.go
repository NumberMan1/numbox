package idcard

import (
	"testing"
)

func TestComprehensive(t *testing.T) {
	testCases := []struct {
		id          string
		description string
	}{
		{"310104196712322X", "格式是否正确: 否 (长度不符)"},
		{"350521199003074452", "大陆居民身份证18位"},
		{"11010119900230123X", "出生日期格式是否有效: 否 (1990-02-30 不是有效日期)"},
		{"11010119900307803X", "大陆居民身份证末位是X18位"},
		{"610104620927690", "大陆居民身份证15位"},
		{"810000199408230021", "港澳居民居住证18位"},
		{"830000199201300022", "台湾居民居住证18位"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			parser := NewIDCardParser(tc.id)

			// 验证有效性
			isValid := parser.IsValid()
			t.Logf("ID: %s, 是否有效: %v", tc.id, isValid)

			if isValid {
				// 获取详细信息
				info, err := parser.GetInfo()
				if err != nil {
					t.Logf("获取信息错误: %v", err)
				} else {
					t.Logf("地区: %s", info.Address)
					t.Logf("出生日期: %s", info.Birthday.Format("2006-01-02"))
					t.Logf("性别: %v", info.Gender)
					t.Logf("星座: %s", info.Constellation)
					t.Logf("生肖: %s", info.ChineseZodiac)
				}

				// 获取性别
				gender, err := parser.GetGender()
				if err != nil {
					t.Logf("获取性别错误: %v", err)
				} else {
					t.Logf("性别: %v", gender)
				}

				// 获取年龄
				age, err := parser.GetAge()
				if err != nil {
					t.Logf("获取年龄错误: %v", err)
				} else {
					t.Logf("年龄: %d", age)
				}
			}
		})
	}
}
