package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func parseAndValidate(t *testing.T, path string) {
	t.Helper()

	content, err := os.ReadFile(path)
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestValidDomainTransferRequest(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_domain_transfer_request.xml")
}

func TestValidDomainTransferQuery(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_domain_transfer_query.xml")
}

func TestValidDomainTransferCancel(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_domain_transfer_cancel.xml")
}

func TestValidDomainTransferApprove(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_domain_transfer_approve.xml")
}

func TestValidDomainTransferReject(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_domain_transfer_reject.xml")
}

func TestValidContactTransferRequest(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_contact_transfer_request.xml")
}

func TestValidContactTransferQuery(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_contact_transfer_query.xml")
}

func TestValidContactTransferCancel(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_contact_transfer_cancel.xml")
}

func TestValidContactTransferApprove(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_contact_transfer_approve.xml")
}

func TestValidContactTransferReject(t *testing.T) {
	parseAndValidate(t, "testdata/transfer/valid_contact_transfer_reject.xml")
}

func TestInvalidTransferOpEmpty(t *testing.T) {
	content, err := os.ReadFile("testdata/transfer/invalid_op_empty.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)
	require.Error(t, err)
}

func TestInvalidTransferOpUnsupported(t *testing.T) {
	content, err := os.ReadFile("testdata/transfer/invalid_op_unsupported.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
}

func TestInvalidTransferTwoObjects(t *testing.T) {
	content, err := os.ReadFile("testdata/transfer/invalid_two_objects.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
}

func TestInvalidTransferNoObject(t *testing.T) {
	content, err := os.ReadFile("testdata/transfer/invalid_no_object.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
}
