package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"io"
	"math/big"
)

var (
	ERR_UNMARSHAL_PUBLIC_KEY = func(err error) error {
		return errors.New("Public key didn't unmarshal: " + err.Error())
	}
	ERR_PUBLIC_KEY_HEX = func(err error) error {
		return errors.New("The public key is not correct hex: " + err.Error())
	}

	ERR_SIGNATURE_HEX = func(err error) error {
		return errors.New("The signature is not correct hex: " + err.Error())
	}

	ERR_SIGNATURE = func(err error) error {
		return errors.New("Could not create signature: " + err.Error())
	}
)

func GenerateKeyPair() (*ecdsa.PrivateKey, string) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)

	publicKeyHex, _ := MarshalPublicKey(&privateKey.PublicKey)
	return privateKey, publicKeyHex
}

func MarshalPublicKey(pub *ecdsa.PublicKey) (string, error) {
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func UnmarshalPublicKey(data string) (*ecdsa.PublicKey, error) {
	b, err := hex.DecodeString(data)
	if err != nil {
		return nil, ERR_PUBLIC_KEY_HEX(err)
	}
	pub, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, ERR_UNMARSHAL_PUBLIC_KEY(err)
	}
	return pub.(*ecdsa.PublicKey), nil
}

func GenerateSharedSecret(privKey *ecdsa.PrivateKey, pubKeyBytes string) ([]byte, error) {
	pubKey, err := UnmarshalPublicKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	x, _ := elliptic.Curve.ScalarMult(elliptic.P521(), pubKey.X, pubKey.Y, privKey.D.Bytes())
	secret := sha256.Sum256(x.Bytes())
	return secret[:], nil
}

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

type Signature struct {
	R, S *big.Int
}

func CreateSignature(privKey *ecdsa.PrivateKey, msg []byte) (string, error) {
	digest := sha256.Sum256(msg)

	r, s, err := ecdsa.Sign(rand.Reader, privKey, digest[:])
	if err != nil {
		return "", ERR_SIGNATURE(err)
	}

	signature := Signature{}
	signature.R = r
	signature.S = s

	var gb bytes.Buffer
	genc := gob.NewEncoder(&gb)
	err = genc.Encode(signature)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(gb.Bytes()), nil
}

func Verify(pubKeyHex, signatureHex string, msg []byte) (bool, error) {
	pubKey, err := UnmarshalPublicKey(pubKeyHex)
	if err != nil {
		return false, err
	}
	signature, err := hex.DecodeString(signatureHex)
	if err != nil {
		return false, ERR_SIGNATURE_HEX(err)
	}
	digest := sha256.Sum256(msg)
	s := new(Signature)
	b := bytes.NewBuffer(signature)
	gdec := gob.NewDecoder(b)
	err = gdec.Decode(s)
	if err != nil {
		return false, err
	}

	return ecdsa.Verify(pubKey, digest[:], s.R, s.S), nil
}
