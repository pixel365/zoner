package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPositiveIntEnv(t *testing.T) {
	type args struct {
		key string
		def int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty", args{"", 1}, 1},
		{"less than zero", args{"TEST_KEY", 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPositiveIntFromEnv(tt.args.key, tt.args.def)
			assert.Equal(t, tt.want, got)
		})
	}
}
