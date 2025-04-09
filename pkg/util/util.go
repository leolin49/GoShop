package util

import (
	"errors"
	"strconv"
	"strings"
)

var Base62Array = [62]byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
	'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
	'U', 'V', 'W', 'X', 'Y', 'Z',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j',
	'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't',
	'u', 'v', 'w', 'x', 'y', 'z',
}

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

func BaseN(x, n int) (s string, err error) {
	if n < 2 || n > 62 {
		err := errors.New("Convert failed from x to n bit, n should in [2, 62]")
		return "", err
	}
	if x < 0 {
		s = "-"
	}
	for ; x > 0; x /= n {
		b := x % n
		s += string(Base62Array[b])
	}
	return
}

func Base62(x int) string {
	s, _ := BaseN(x, 62)
	return s
}
