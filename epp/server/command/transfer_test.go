package command

import (
	"testing"
)

func TestValidDomainTransferRequest(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_domain_transfer_request.xml")
}

func TestValidDomainTransferQuery(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_domain_transfer_query.xml")
}

func TestValidDomainTransferCancel(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_domain_transfer_cancel.xml")
}

func TestValidDomainTransferApprove(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_domain_transfer_approve.xml")
}

func TestValidDomainTransferReject(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_domain_transfer_reject.xml")
}

func TestValidContactTransferRequest(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_contact_transfer_request.xml")
}

func TestValidContactTransferQuery(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_contact_transfer_query.xml")
}

func TestValidContactTransferCancel(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_contact_transfer_cancel.xml")
}

func TestValidContactTransferApprove(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_contact_transfer_approve.xml")
}

func TestValidContactTransferReject(t *testing.T) {
	mustParseAndValidate(t, "testdata/transfer/valid_contact_transfer_reject.xml")
}

func TestInvalidTransferOpEmpty(t *testing.T) {
	mustFailParse(t, "testdata/transfer/invalid_op_empty.xml")
}

func TestInvalidTransferOpUnsupported(t *testing.T) {
	mustFailParse(t, "testdata/transfer/invalid_op_unsupported.xml")
}

func TestInvalidTransferTwoObjects(t *testing.T) {
	mustFailParse(t, "testdata/transfer/invalid_two_objects.xml")
}

func TestInvalidTransferNoObject(t *testing.T) {
	mustFailParse(t, "testdata/transfer/invalid_no_object.xml")
}
