package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	command2 "github.com/pixel365/zoner/epp/server/command/command"
)

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	return b
}

func mustParse(t *testing.T, path string) command2.Commander {
	t.Helper()
	p := CmdParser{}
	cmd, err := p.Parse(mustRead(t, path))
	require.NoError(t, err)
	return cmd
}

func mustParseAndValidate(t *testing.T, path string) {
	t.Helper()
	cmd := mustParse(t, path)
	require.NoError(t, cmd.Validate())
}

func mustFailParse(t *testing.T, path string) {
	t.Helper()
	p := CmdParser{}
	_, err := p.Parse(mustRead(t, path))
	require.Error(t, err)
}
