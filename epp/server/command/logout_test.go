package command

import (
	"testing"
)

func TestValidLogout(t *testing.T) {
	mustParseAndValidate(t, "testdata/logout/valid.xml")
}

func TestInvalidLogout(t *testing.T) {
	mustFailParse(t, "testdata/logout/invalid.xml")
}
