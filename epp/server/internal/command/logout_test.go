package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidLogout(t *testing.T) {
	content, err := os.ReadFile("testdata/logout/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestInvalidLogout(t *testing.T) {
	content, err := os.ReadFile("testdata/logout/invalid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)

	assert.Contains(t, err.Error(), "<logout> must be empty")
}
