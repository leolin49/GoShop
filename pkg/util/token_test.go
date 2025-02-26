package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJwtToken(t *testing.T) {
	userId := uint32(49)
	accessSec := int64(3600)
	refreshSec := int64(86400)

	accessToken, refreshToken, err := JwtDoubleToken(userId, accessSec, refreshSec)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	accessUserId, err := JwtExtractAccessTokenUserId(accessToken)
	assert.NoError(t, err)
	assert.Equal(t, userId, accessUserId)

	refreshUserId, err := JwtExtractRefreshTokenUserId(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, userId, refreshUserId)
}
