package util

func StrEq(a, b interface{}) bool {
	c, ok := a.(string)
	if !ok {
		return false
	}
	d, ok := b.(string)
	if !ok {
		return false
	}
	if c == d {
		return true
	}
	return false
}
