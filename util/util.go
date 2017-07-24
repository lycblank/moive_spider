package util

// Min求最小值
func Min(left, right interface{}) interface{} {
	switch left.(type) {
	case uint32:
		l, _ := left.(uint32)
		r, _ := right.(uint32)
		if l < r {
			return l
		}
		return r
	}
	return left
}
