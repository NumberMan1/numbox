package typeutil

// Integer includes signed and unsigned integers and their named types.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Number includes types that support arithmetic. Map keys only need comparable.
type Number interface {
	Integer | ~float32 | ~float64
}
