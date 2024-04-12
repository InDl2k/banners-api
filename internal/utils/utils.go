package utils

import "strconv"

func ParseOrDefaultInt(val string, def int) int {
	res, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return res
}

func ParseOrDefaultBool(val string, def bool) bool {
	res, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return res
}
