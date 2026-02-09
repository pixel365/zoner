package command

import (
	"testing"
)

func TestValidOpReq(t *testing.T) {
	mustParseAndValidate(t, "testdata/poll/valid_req.xml")
}

func TestValidOpAck(t *testing.T) {
	mustParseAndValidate(t, "testdata/poll/valid_ack.xml")
}

func TestEmptyOp(t *testing.T) {
	mustFailParse(t, "testdata/poll/empty_op.xml")
}

func TestInvalidOp(t *testing.T) {
	mustFailParse(t, "testdata/poll/invalid_op.xml")
}
