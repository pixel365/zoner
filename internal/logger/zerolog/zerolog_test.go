package zerolog

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/pixel365/zoner/internal/logger"
	"github.com/stretchr/testify/assert"
)

type w struct{}

func (o *w) Write(p []byte) (n int, err error) {
	f, _ := os.OpenFile("test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return f.Write(p)
}

func TestNewLogger(t *testing.T) {
	t.Parallel()

	f, err := os.Create("test.log")
	if err != nil {
		t.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}

		_ = os.Remove("test.log")
	}(f)

	writers := make([]io.Writer, 0, 1)
	writers = append(writers, &w{})

	l := NewLogger(logger.NewConfig(), writers...)

	assert.NotNil(t, l)

	l.Debug("a")
	l.Info("b")
	l.Error("c", errors.New("c"))

	i := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i++
		switch i {
		case 1:
			assert.True(t, strings.HasSuffix(scanner.Text(), "\"message\":\"b\"}"))
		case 2:
			assert.True(t, strings.HasSuffix(scanner.Text(), "\"message\":\"c\"}"))
		}
	}

	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 2, i)
}
