package command

import (
	"testing"
)

func TestValidDomainUpdateAddNS(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_add_ns.xml")
}

func TestValidDomainUpdateRemContact(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_rem_contact.xml")
}

func TestValidDomainUpdateAddStatus(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_add_status.xml")
}

func TestValidDomainUpdateChgRegistrant(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_chg_registrant.xml")
}

func TestValidDomainUpdateChgAuthInfo(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_chg_authinfo.xml")
}

func TestValidDomainUpdateAll(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_add_rem_chg.xml")
}

func TestValidDomainUpdateChgRemoveRegistrant(t *testing.T) {
	mustParseAndValidate(t, "testdata/update/valid_chg_remove_registrant.xml")
}

func TestInvalidDomainUpdateNoObject(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_no_object.xml")
}

func TestInvalidDomainUpdateWrongNamespace(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_wrong_namespace.xml")
}

func TestInvalidDomainUpdateTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_two_objects.xml")
}

func TestInvalidDomainUpdateMissingName(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_missing_name.xml")
}

func TestInvalidDomainUpdateNoAddRemChg(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_no_add_rem_chg.xml")
}

func TestInvalidDomainUpdateAddEmpty(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_add_empty.xml")
}

func TestInvalidDomainUpdateRemEmpty(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_rem_empty.xml")
}

func TestInvalidDomainUpdateChgEmpty(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_chg_empty.xml")
}

func TestValidDomainUpdateChgAuthIbfoEmpty(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_chg_authinfo_empty.xml")
}

func TestValidDomainUpdateChgAuthInfoAndNull(t *testing.T) {
	mustFailParse(t, "testdata/update/invalid_chg_authinfo_pw_and_null.xml")
}
