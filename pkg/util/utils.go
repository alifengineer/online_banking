package util

import (
	"strconv"
	"time"
)

// IfElse evaluates a condition, if true returns the first parameter otherwise the second
func IfElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IsValidTimeStamp(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
