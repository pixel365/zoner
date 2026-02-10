package command

import (
	"testing"
)

func TestValidDomainDelete(t *testing.T) {
	mustParseAndValidate(t, "testdata/delete/valid_domain_delete.xml")
}

func TestValidContactDelete(t *testing.T) {
	mustParseAndValidate(t, "testdata/delete/valid_contact_delete.xml")
}

func TestValidHostDelete(t *testing.T) {
	mustParseAndValidate(t, "testdata/delete/valid_host_delete.xml")
}

func TestInvalidDeleteNoObject(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_no_object.xml")
}

func TestInvalidDeleteTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_two_objects.xml")
}

func TestInvalidDeleteUnsupportedNamespace(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_unsupported_namespace.xml")
}

func TestInvalidDomainDeleteMissingName(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_domain_missing_name.xml")
}

func TestInvalidContactDeleteMissingID(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_contact_missing_id.xml")
}

func TestInvalidHostDeleteMissingName(t *testing.T) {
	mustFailParse(t, "testdata/delete/invalid_host_missing_name.xml")
}
