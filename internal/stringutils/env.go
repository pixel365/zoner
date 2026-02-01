package stringutils

import (
	"os"
	"strconv"
)

func GetPositiveIntFromEnv(key string, def int) int {
	v, _ := strconv.Atoi(os.Getenv(key))
	if v < 1 {
		return def
	}
	return v
}
