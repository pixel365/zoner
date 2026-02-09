package command

import (
	"testing"
)

func TestValidHostCreateMin(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/host/valid_min.xml")
}

func TestValidHostCreateWithAddrs(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/host/valid_with_addrs.xml")
}

func TestInvalidHostCreateMissingName(t *testing.T) {
	mustFailParse(t, "testdata/create/host/invalid_missing_name.xml")
}

func TestInvalidHostCreateBadIPAttr(t *testing.T) {
	mustFailParse(t, "testdata/create/host/invalid_addr_bad_ip_attr.xml")
}

func TestValidDomainCreateMin(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/domain/valid_min.xml")
}

func TestValidDomainCreateFullHostObj(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/domain/valid_full_hostObj.xml")
}

func TestValidDomainCreateNSHostAttr(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/domain/valid_ns_hostAttr.xml")
}

func TestInvalidDomainCreateMissingName(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_missing_name.xml")
}

func TestInvalidDomainCreateMissingAuthInfo(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_missing_authinfo.xml")
}

func TestInvalidDomainCreatePeriodBadUnit(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_period_bad_unit.xml")
}

func TestInvalidDomainCreateNsEmpty(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_ns_empty.xml")
}

func TestInvalidDomainCreateHostAttrMissingHostName(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_ns_hostAttr_missing_hostName.xml")
}

func TestValidContactCreateMin(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/contact/valid_min.xml")
}

func TestValidContactCreateFull(t *testing.T) {
	mustParseAndValidate(t, "testdata/create/contact/valid_full.xml")
}

func TestInvalidContactCreateMissingID(t *testing.T) {
	mustFailParse(t, "testdata/create/domain/invalid_period_bad_unit.xml")
}

func TestInvalidContactCreateMissingPostalInfo(t *testing.T) {
	mustFailParse(t, "testdata/create/contact/invalid_missing_postalInfo.xml")
}

func TestInvalidContactCreatePostalInfoBadType(t *testing.T) {
	mustFailParse(t, "testdata/create/contact/invalid_postalInfo_bad_type.xml")
}

func TestInvalidContactCreateAddrMissingCity(t *testing.T) {
	mustFailParse(t, "testdata/create/contact/invalid_addr_missing_city.xml")
}

func TestInvalidContactCreateMissingEmail(t *testing.T) {
	mustFailParse(t, "testdata/create/contact/invalid_missing_email.xml")
}

func TestInvalidContactCreateMissingAuthInfo(t *testing.T) {
	mustFailParse(t, "testdata/create/contact/invalid_missing_authinfo.xml")
}

func TestInvalidCreateNoObject(t *testing.T) {
	mustFailParse(t, "testdata/create/invalid_no_object.xml")
}

func TestInvalidCreateTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/create/invalid_two_objects.xml")
}
