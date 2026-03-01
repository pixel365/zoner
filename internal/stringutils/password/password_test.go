package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndVerify(t *testing.T) {
	t.Parallel()

	raw := "S3cure-P@ssw0rd"

	encoded, err := Hash(raw, DefaultParams)
	require.NoError(t, err)

	ok, err := Verify(raw, encoded)
	require.NoError(t, err)
	assert.True(t, ok)
}

func TestVerifyWrongPassword(t *testing.T) {
	t.Parallel()

	encoded, err := Hash("correct-password", DefaultParams)
	require.NoError(t, err)

	ok, err := Verify("wrong-password", encoded)
	require.NoError(t, err)
	assert.False(t, ok)
}

func TestVerifyInvalidFormat(t *testing.T) {
	t.Parallel()

	ok, err := Verify("anything", "not-a-valid-argon2-hash")
	require.Error(t, err)
	assert.False(t, ok)
}
