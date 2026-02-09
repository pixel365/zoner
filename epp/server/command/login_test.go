package command

import (
	"testing"
)

func TestLogin(t *testing.T) {
	mustParseAndValidate(t, "testdata/login/valid.xml")
}
