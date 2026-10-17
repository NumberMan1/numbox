# randutil

RandomInt(seed,min,max) 使用 [min,max)，相等时返回 min；RandInt(min,max) 使用 [min,max]。
逆序区间会 panic。调用方需保证结果能装入返回类型。

Shuffle 使用显式种子，ShuffleRandom 使用时间种子；两者就地修改切片。
保留原随机算法与相同种子下的行为。随机数不用于安全令牌或密码学。
