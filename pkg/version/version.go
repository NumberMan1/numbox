package version

import (
	"math"
	"strconv"
	"strings"
)

type Version string

func (ver Version) String() string {
	return string(ver)
}

func (ver Version) PkgSuffix() string {
	var versionStr, weekVerStr = "", "0"
	verArr := strings.Split(ver.TrimPrefix().TrimSuffix().String(), ".")
	if len(ver) >= 4 {
		weekVerStr = verArr[3]
	}
	versionStr = strings.Join(verArr[:3], ".") + "-" + weekVerStr
	return versionStr
}

func (ver Version) Tag() string {
	var versionTag, weekVerStr string
	verArr := strings.Split(ver.TrimSuffix().String(), ".")
	if len(ver) >= 3 {
		weekVerStr = verArr[3]
	}
	verArr = verArr[:3]
	versionTag = strings.Join(verArr, ".")
	if len(weekVerStr) > 0 {
		versionTag += "-" + weekVerStr
	}
	versionTag = strings.TrimPrefix(versionTag, "v")
	return versionTag
}

func (ver Version) TrimSuffix() Version {
	verStr := ver.String()
	idx := strings.LastIndex(verStr, "-") // 找到最后一个 '-' 的位置
	if idx != -1 {
		verStr = verStr[:idx] // 截取 '-' 之前的部分
	}
	return Version(verStr)
}

func (ver Version) TrimPrefix() Version {
	newVer := strings.TrimPrefix(ver.String(), "v")
	return Version(newVer)
}

func (ver Version) CompatibleVersion() Version {
	var compatibleVer = ""
	verArr := strings.Split(ver.TrimSuffix().String(), ".")
	if len(verArr) > 3 {
		verArr = []string{verArr[0], verArr[1], "0", verArr[3]}
	} else {
		verArr = verArr[:2]
		verArr = append(verArr, "0")
	}
	compatibleVer = strings.Join(verArr, ".")
	return Version(compatibleVer)
}

func (ver Version) Number() int64 {
	verStr := ver.TrimSuffix().String()
	verStr = strings.TrimPrefix(verStr, "v")
	verArr := strings.Split(verStr, ".")
	if len(verArr) > 4 {
		verArr = verArr[:4]
	}

	var verNum int64
	for i, str := range verArr {
		if i < 3 {
			v, _ := strconv.ParseInt(str, 10, 64)
			verNum += v * int64(math.Pow(1000, float64(4-i)))
		} else {
			v, _ := strconv.ParseInt(str, 10, 64)
			verNum += v
		}
	}
	return verNum
}

func (ver Version) Compatibility(otherVer Version) bool {
	verNum := ver.CompatibleVersion().Number()
	otherVerNum := otherVer.CompatibleVersion().Number()
	if verNum/1000000 == otherVerNum/1000000 && verNum >= otherVerNum {
		return true
	}
	return false
}

func (ver Version) Equal(otherVer Version) bool {
	return ver.String() == otherVer.String()
}

func (ver Version) IsLowerThan(otherVer Version) bool {
	return ver.Number() < otherVer.Number()
}
