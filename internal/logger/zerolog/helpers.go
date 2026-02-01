package zerolog

import (
	"os"
	"runtime"
)

func hostname() string {
	if v := os.Getenv("HOSTNAME"); v != "" {
		return v
	}

	if h, _ := os.Hostname(); h != "" {
		return h
	}

	return runtime.GOOS
}
