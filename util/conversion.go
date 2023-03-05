package util

import "strconv"

func ParseFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func FormatScientific(f float64) string {
	return strconv.FormatFloat(f, 'e', 16, 64)
}
