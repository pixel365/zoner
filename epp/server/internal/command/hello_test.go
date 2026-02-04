package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidHello(t *testing.T) {
	content, err := os.ReadFile("testdata/hello/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestInvalidHello(t *testing.T) {
	content, err := os.ReadFile("testdata/hello/invalid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)

	assert.Contains(t, err.Error(), "<hello> must be empty")
}
