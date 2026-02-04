package command

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidDomainInfo(t *testing.T) {
	content, err := os.ReadFile("testdata/info/domain/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestValidDomainInfoWithAuth(t *testing.T) {
	content, err := os.ReadFile("testdata/info/domain/valid_auth.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestValidDomainInfoWithEmptyPassword(t *testing.T) {
	content, err := os.ReadFile("testdata/info/domain/invalid_auth.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
}

func TestValidContactInfo(t *testing.T) {
	content, err := os.ReadFile("testdata/info/contact/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestValidContactInfoWithAuth(t *testing.T) {
	content, err := os.ReadFile("testdata/info/contact/valid_auth.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}

func TestInvalidContactInfoWithEmptyPassword(t *testing.T) {
	content, err := os.ReadFile("testdata/info/contact/invalid_auth.xml")
	require.NoError(t, err)

	p := CmdParser{}
	_, err = p.Parse(content)

	require.Error(t, err)
}

func TestValidHostInfo(t *testing.T) {
	content, err := os.ReadFile("testdata/info/host/valid.xml")
	require.NoError(t, err)

	p := CmdParser{}
	cmd, err := p.Parse(content)

	require.NoError(t, err)
	require.NoError(t, cmd.Validate())
}
