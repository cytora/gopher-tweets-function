package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeDecode(t *testing.T) {
	key := "adelina-123-kgfglhkfmglhfmghlfg"
	creds := &Credentials{
		ExtAccessKey:    "kinda-secret",
		ExtAccessSecret: "super-secret",
	}
	encoded := creds.Encode(key)
	empty := creds.Encode("")
	assert.NotEqual(t, empty, encoded)
	assert.NotEmpty(t, encoded)
	decodedCreds, err := Decode(key, encoded)
	assert.Nil(t, err)
	assert.NotNil(t, decodedCreds)
	assert.Equal(t, *creds, *decodedCreds)
	fail, err := Decode("wrong", encoded)
	assert.Nil(t, fail)
}
