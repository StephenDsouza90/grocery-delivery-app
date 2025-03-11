package utils

import (
	"strconv"

	"github.com/google/uuid"
)

// ConverterStrToInt converts a string to an integer
func ConverterStrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

// ConverterIntToStr converts an integer to a string
func ConverterIntToStr(i int) string {
	return strconv.Itoa(i)
}

// GenerateUUID generates unique a UUID
func GenerateUUID() string {
	return uuid.New().String()

}
