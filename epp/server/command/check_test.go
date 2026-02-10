package command

import (
	"testing"
)

func TestValidDomainCheck(t *testing.T) {
	mustParseAndValidate(t, "testdata/check/valid_domain_check.xml")
}

func TestValidContactCheck(t *testing.T) {
	mustParseAndValidate(t, "testdata/check/valid_contact_check.xml")
}

func TestValidHostCheck(t *testing.T) {
	mustParseAndValidate(t, "testdata/check/valid_host_check.xml")
}

func TestInvalidCheckNoObject(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_no_object.xml")
}

func TestInvalidCheckTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_two_objects.xml")
}

func TestInvalidCheckUnsupportedNamespace(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_unsupported_namespace.xml")
}

func TestInvalidDomainCheckEmptyList(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_domain_empty_list.xml")
}

func TestInvalidContactCheckEmptyList(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_contact_empty_list.xml")
}

func TestInvalidHostCheckEmptyList(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_host_empty_list.xml")
}

func TestInvalidDomainCheckHasEmptyName(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_domain_has_empty_name.xml")
}

func TestInvalidContactCheckHasEmptyID(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_contact_has_empty_id.xml")
}

func TestInvalidHostCheckHasEmptyName(t *testing.T) {
	mustFailParse(t, "testdata/check/invalid_host_has_empty_name.xml")
}
