package percent

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func ParsePercentageToPermille(input string) (int, error) {
	// 去掉空格并检查是否包含百分号
	input = strings.TrimSpace(input)
	if !strings.HasSuffix(input, "%") {
		return 0, fmt.Errorf("输入格式错误，必须包含百分号")
	}

	// 去掉百分号并解析浮点数
	numStr := strings.TrimSuffix(input, "%")
	value, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, err
	}

	// 计算千分比并四舍五入
	permille := math.Round(value * 10)
	if math.IsNaN(permille) || math.IsInf(permille, 0) || permille >= float64(int(^uint(0)>>1)) || permille < -float64(int(^uint(0)>>1))-1 {
		return 0, fmt.Errorf("percentage outside integer range")
	}
	return int(permille), nil
}
