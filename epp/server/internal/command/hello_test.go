package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	content, err := os.ReadFile("testdata/hello.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}
