package utils

import (
	"strconv"
)

// string to int
func Str2Int(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}
	return i
}

// string to int64
func Str2Int64(str string) int64 {
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// int to string
func Int2Str(i int) string {
	return strconv.Itoa(i)
}

// int64 to string
func Int642Str(i int64) string {
	return strconv.FormatInt(i, 10)
}
