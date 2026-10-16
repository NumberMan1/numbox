package idcard

import (
	"github.com/guanguans/id-validator"
	"time"
)

type Gender int

const (
	GenderUnknown Gender = 0
	GenderFemale  Gender = 1
	GenderMan     Gender = 2
)

type IdInfo struct {
	AddressCode   int
	Abandoned     int
	Address       string
	AddressTree   []string
	Birthday      time.Time
	Constellation string
	ChineseZodiac string
	Gender        Gender
	Length        int
	CheckBit      string
}

// IDCardParser 身份证解析器结构体
type IDCardParser struct {
	id string
}

// NewIDCardParser 创建新的身份证解析器
func NewIDCardParser(idStr string) *IDCardParser {
	return &IDCardParser{
		id: idStr,
	}
}

// IsValid 验证身份证号是否有效
func (p *IDCardParser) IsValid() bool {
	return idvalidator.IsValid(p.id, false)
}

// GetBirthDate 获取出生日期
func (p *IDCardParser) GetBirthDate() (time.Time, error) {
	info, err := idvalidator.GetInfo(p.id, false)
	if err != nil {
		return time.Time{}, err
	}

	return info.Birthday, nil
}

// GetGender 获取性别
func (p *IDCardParser) GetGender() (Gender, error) {
	info, err := idvalidator.GetInfo(p.id, false)
	if err != nil {
		return GenderUnknown, err
	}

	if info.Sex == 1 {
		return GenderMan, nil
	}
	return GenderFemale, nil
}

// GetAge 获取年龄
func (p *IDCardParser) GetAge() (int, error) {
	info, err := idvalidator.GetInfo(p.id, false)
	if err != nil {
		return 0, err
	}

	now := time.Now()
	age := now.Year() - info.Birthday.Year()

	// 如果还没到生日，年龄减1
	if now.Month() < info.Birthday.Month() ||
		(now.Month() == info.Birthday.Month() && now.Day() < info.Birthday.Day()) {
		age--
	}

	return age, nil
}

// GetInfo 获取身份证信息
func (p *IDCardParser) GetInfo() (IdInfo, error) {
	info, err := idvalidator.GetInfo(p.id, false)
	if err != nil {
		return IdInfo{}, err
	}
	idInfo := IdInfo{
		AddressCode:   info.AddressCode,
		Abandoned:     info.Abandoned,
		Address:       info.Address,
		AddressTree:   info.AddressTree,
		Birthday:      info.Birthday,
		Constellation: info.Constellation,
		ChineseZodiac: info.ChineseZodiac,
		Length:        info.Length,
		CheckBit:      info.CheckBit,
	}
	if info.Sex == 1 {
		idInfo.Gender = GenderMan
	} else {
		idInfo.Gender = GenderFemale
	}
	return idInfo, nil
}
