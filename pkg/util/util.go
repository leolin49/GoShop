package util

import (
	"strconv"
	"strings"
)

func DecToBin(x int) string {
	return strconv.FormatInt(int64(x), 2)
}

func BinToDec(s string) (int64, error) {
	return strconv.ParseInt(s, 2, 64)
}

func DecToOct(x int) string {
	return strconv.FormatInt(int64(x), 8)
}

func OctToDec(s string) (int64, error) {
	return strconv.ParseInt(s, 8, 64)
}

func DecToHex(x int) string {
	return strings.ToUpper(strconv.FormatUint(uint64(x), 16))
}

func HexToDec(s string) (int64, error) {
	return strconv.ParseInt(s, 16, 64)
}
