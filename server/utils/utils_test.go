package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcdh(t *testing.T) {
	priv1, pub1 := GenerateKeyPair()
	priv2, pub2 := GenerateKeyPair()

	sh1, err := GenerateSharedSecret(priv1, pub2)
	assert.Nil(t, err)
	sh2, err := GenerateSharedSecret(priv2, pub1)
	assert.Nil(t, err)
	assert.Equal(t, sh1, sh2)
}

func TestEncryption(t *testing.T) {
	priv1, _ := GenerateKeyPair()
	_, pub2 := GenerateKeyPair()
	sh, err := GenerateSharedSecret(priv1, pub2)
	assert.Nil(t, err)

	msg := []byte("lalalallala")

	encrypted, err := Encrypt(msg, sh)
	assert.Nil(t, err)

	decrypted, err := Decrypt(encrypted, sh)
	assert.Nil(t, err)

	assert.Equal(t, msg, decrypted)
}

func TestSign(t *testing.T) {
	priv, pub := GenerateKeyPair()
	msg := []byte("lalalalalala")
	signature, err := CreateSignature(priv, msg)
	assert.Nil(t, err)
	v, err := Verify(pub, signature, msg)
	assert.Nil(t, err)
	assert.True(t, v)
}
