package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidCheck(t *testing.T) {
	content, err := os.ReadFile("testdata/check/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestInvalidCheck(t *testing.T) {
	content, err := os.ReadFile("testdata/check/invalid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
	assert.Equal(t, "unmarshal xml payload: objects is empty", err.Error())
}
