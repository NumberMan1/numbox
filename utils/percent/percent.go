package percent

import "math/rand/v2"

const (
	FullPercent Percent = 1000
	HalfPercent Percent = 500
)

type Percent int32

func (per Percent) Int() int {
	return int(per)
}

func (per Percent) Int64() int64 {
	return int64(per)
}

func (per Percent) Int32() int32 {
	return int32(per)
}

func (per Percent) Float64() float64 {
	return float64(per) / float64(FullPercent)
}

func (per Percent) Hit() bool {
	randRes := rand.IntN(FullPercent.Int())
	return randRes <= per.Int()
}

func (per Percent) HitBySeed(seed int64) bool {
	r := rand.New(rand.NewPCG(uint64(seed), uint64(seed)))
	randRes := r.IntN(FullPercent.Int())
	return randRes <= per.Int()
}

func (per Percent) Multiply(val Percent) Percent {
	return Percent(per.Float64() * val.Float64() * 1000)
}

func (per Percent) MultiplyByInt(val int) int {
	return int(per.Float64() * float64(val))
}

func (per Percent) MultiplyByInt32(val int32) int32 {
	return int32(per.Float64() * float64(val))
}

func (per Percent) MultiplyByInt64(val int64) int64 {
	return int64(per.Float64() * float64(val))
}

func GetPercentByDivide[T int | int8 | int32 | int64](dividend, divisor T) Percent {
	result := float64(dividend) / float64(divisor)
	result = result * 1000
	return Percent(result)
}
