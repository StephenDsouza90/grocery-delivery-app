package utils

import "strconv"

func ConverterStrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ConverterIntToStr(i int) string {
	return strconv.Itoa(i)
}
