package command

import (
	"testing"
)

func TestValidDomainRenewMin(t *testing.T) {
	mustParseAndValidate(t, "testdata/renew/valid_min.xml")
}

func TestValidDomainRenewWithPeriod(t *testing.T) {
	mustParseAndValidate(t, "testdata/renew/valid_with_period.xml")
}

func TestInvalidRenewNoObject(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_no_object.xml")
}

func TestInvalidRenewWrongNamespace(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_wrong_namespace.xml")
}

func TestInvalidRenewTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_two_objects.xml")
}

func TestInvalidRenewMissingName(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_missing_name.xml")
}

func TestInvalidRenewMissingCurExpDate(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_missing_curExpDate.xml")
}

func TestInvalidRenewPeriodBadUnit(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_period_bad_unit.xml")
}

func TestInvalidRenewPeriodNonPositive(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_period_non_positive.xml")
}

func TestInvalidRenewCurExpDateBadFormat(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_curExpDate_bad_format.xml")
}

func TestInvalidRenewCurExpDateImpossible(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_curExpDate_impossible.xml")
}

func TestInvalidRenewCurExpDateHasTime(t *testing.T) {
	mustFailParse(t, "testdata/renew/invalid_curExpDate_has_time.xml")
}
