package command

import (
	"testing"
)

func TestValidDomainInfo(t *testing.T) {
	mustParseAndValidate(t, "testdata/info/domain/valid.xml")
}

func TestValidDomainInfoWithAuth(t *testing.T) {
	mustParseAndValidate(t, "testdata/info/domain/valid_auth.xml")
}

func TestValidDomainInfoWithEmptyPassword(t *testing.T) {
	mustFailParse(t, "testdata/info/domain/invalid_auth.xml")
}

func TestValidContactInfo(t *testing.T) {
	mustParseAndValidate(t, "testdata/info/contact/valid.xml")
}

func TestValidContactInfoWithAuth(t *testing.T) {
	mustParseAndValidate(t, "testdata/info/contact/valid_auth.xml")
}

func TestInvalidContactInfoWithEmptyPassword(t *testing.T) {
	mustFailParse(t, "testdata/info/contact/invalid_auth.xml")
}

func TestValidHostInfo(t *testing.T) {
	mustParseAndValidate(t, "testdata/info/host/valid.xml")
}
