package command

import (
	"testing"
)

func TestValidCheck(t *testing.T) {
	mustParseAndValidate(t, "testdata/check/valid.xml")
}

func TestInvalidCheck(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid.xml")
}
