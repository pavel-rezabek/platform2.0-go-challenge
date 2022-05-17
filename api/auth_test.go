package api

import (
	"testing"
)

// TestGenerateJWT calls api.GenerateJWT with and id checking,
// that it produces valid JWT token
func TestGenerateJWT(t *testing.T) {
	id := uint(0)
	token, err := GenerateJWT(id)
	_, parseErr := ParseToken(token)
	if parseErr != nil || err != nil {
		t.Fatalf(`GenerateJWT(uint(0)) = %q, %v, parseErr: %v`, token, err, parseErr)
	}
}
