package percent

import (
	"fmt"
	"testing"
	"time"
)

func TestPercent(t *testing.T) {
	var (
		randomSeed                = time.Now().UnixMilli()
		per               Percent = 700
		totalTimes        int64
		totalSucceedTimes int64
	)
	for times := 0; times < 1000; times++ {
		var nodeNo int64
		for nodeNo = 102001; nodeNo < 102018; nodeNo++ {
			for i := 0; i < 3; i++ {
				totalTimes++
				testRandomSeed := randomSeed*1000000 + nodeNo*1000 + int64(i+1)
				if per.HitBySeed(testRandomSeed) {
					totalSucceedTimes++
				}
			}
		}
	}

	fmt.Println("total times:", totalTimes)
	fmt.Println("total succeed times:", totalSucceedTimes)
	fmt.Println("average succeed times:", totalSucceedTimes/1000)
}
