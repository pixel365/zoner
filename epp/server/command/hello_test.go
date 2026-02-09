package command

import (
	"testing"
)

func TestValidHello(t *testing.T) {
	mustParseAndValidate(t, "testdata/hello/valid.xml")
}

func TestInvalidHello(t *testing.T) {
	mustFailParse(t, "testdata/hello/invalid.xml")
}
