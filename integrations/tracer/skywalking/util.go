package skywalking

import "strconv"

// stringToSpanID converts a string to int32.
// If the string given is not a valid integer, it returns 0.
func stringToSpanID(s string) int32 {
	id, _ := strconv.ParseInt(s, 10, 32)
	return int32(id)
}
