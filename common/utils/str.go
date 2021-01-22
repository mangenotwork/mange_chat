package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
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

//随机字符
const char = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandChar(size int) string {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(char[rand.Int63()%int64(len(char))])
	}
	return s.String()
}

//DateTime DateTime
func Unix2Date(unix int64) (ret string) {
	if unix < 1 {
		return
	}
	tm := time.Unix(unix, 0)
	ret = FormatTime(tm, "YYYY-MM-DD HH:mm")
	return
}

//DateTime DateTime
func StrUnix2Date(unixStr string) (ret string) {

	unix := Str2Int64(unixStr)

	if unix < 1 {
		return
	}
	tm := time.Unix(unix, 0)
	ret = FormatTime(tm, "YYYY-MM-DD HH:mm")
	return
}

//DateTimeFormat DateTimeFormat
func DateTimeFormat(unix int64, format string) (ret string) {
	if unix < 1 {
		return
	}
	tm := time.Unix(unix, 0)
	ret = FormatTime(tm, format)
	return
}

//FormatTime 格式化时间显示
func FormatTime(t time.Time, format string) string {
	res := strings.Replace(format, "MM", t.Format("01"), -1)
	res = strings.Replace(res, "M", t.Format("1"), -1)
	res = strings.Replace(res, "DD", t.Format("02"), -1)
	res = strings.Replace(res, "D", t.Format("2"), -1)
	res = strings.Replace(res, "YYYY", t.Format("2006"), -1)
	res = strings.Replace(res, "YY", t.Format("06"), -1)
	res = strings.Replace(res, "HH", fmt.Sprintf("%02d", t.Hour()), -1)
	res = strings.Replace(res, "H", fmt.Sprintf("%d", t.Hour()), -1)
	res = strings.Replace(res, "hh", t.Format("03"), -1)
	res = strings.Replace(res, "h", t.Format("3"), -1)
	res = strings.Replace(res, "mm", t.Format("04"), -1)
	res = strings.Replace(res, "m", t.Format("4"), -1)
	res = strings.Replace(res, "ss", t.Format("05"), -1)
	res = strings.Replace(res, "s", t.Format("5"), -1)
	return res
}
