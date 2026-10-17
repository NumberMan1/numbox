package percent

import (
	"github.com/NumberMan1/numbox/utils/convutil"
)

func FormatPercentString(percent Percent) string {
	return convutil.FormatIntString(percent.Int32())
}

func ParsePercentString(src string) Percent {
	return Percent(convutil.ParseIntString[int32](src))
}
