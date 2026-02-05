package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidOpReq(t *testing.T) {
	content, err := os.ReadFile("testdata/poll/valid_req.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestValidOpAck(t *testing.T) {
	content, err := os.ReadFile("testdata/poll/valid_ack.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestEmptyOp(t *testing.T) {
	content, err := os.ReadFile("testdata/poll/empty_op.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
	assert.Equal(t, "invalid poll operation: ", err.Error())
}

func TestInvalidOp(t *testing.T) {
	content, err := os.ReadFile("testdata/poll/invalid_op.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
	assert.Equal(t, "invalid poll operation: xxx", err.Error())
}
