package ctrls

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mragiadakos/tenderforwarding/server/confs"
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
	"github.com/mragiadakos/tenderforwarding/server/validators"
	uuid "github.com/satori/go.uuid"

	"github.com/stretchr/testify/assert"
)

func TestReceivedFailOnHashDoesNotExists(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	d := models.DeliveryModel{}
	d.Action = models.RECEIVED
	d.Data = models.ReceivedModel{
		Receiver: pubHex,
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_HASH_DOES_NOT_EXISTS, errors.New(resp.Log))
}

func createForward(t *testing.T, app *TFApplication, receiverPubHex string) string {
	priv, pubHex := utils.GenerateKeyPair()

	confs.Conf.Redistributors = []string{pubHex}
	coin := uuid.NewV4().String()

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	fm := models.ForwardModel{
		Redistributor:   pubHex,
		Coins:           []string{coin},
		EncryptedOwners: map[string]string{coin: "lalalla fake encryption"},
		Metadata:        map[string]string{"alal": "lalala"},
		Receiver:        receiverPubHex,
	}
	d.Data = fm
	msg, _ := json.Marshal(d.Data)
	signature, err := utils.CreateSignature(priv, msg)
	assert.Nil(t, err)
	d.Signature = signature

	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeOK, resp.Code)

	cb, _ := json.Marshal(fm.Coins)
	hash := sha256.Sum256(cb)
	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

func TestReceivedFailOnNotReceiver(t *testing.T) {
	app := NewTFApplication()
	_, receiverPubHex := utils.GenerateKeyPair()
	_, fakePubHex := utils.GenerateKeyPair()

	hash := createForward(t, app, receiverPubHex)

	d := models.DeliveryModel{}
	d.Action = models.RECEIVED
	d.Data = models.ReceivedModel{
		Receiver: fakePubHex,
		Hash:     hash,
	}

	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_RECEIVER_NOT_IN_FORWARD, errors.New(resp.Log))
}

func TestReceiverFailOnSignature(t *testing.T) {
	app := NewTFApplication()
	_, receiverPubHex := utils.GenerateKeyPair()
	fakePriv, _ := utils.GenerateKeyPair()

	hash := createForward(t, app, receiverPubHex)

	d := models.DeliveryModel{}
	d.Action = models.RECEIVED
	rm := models.ReceivedModel{
		Receiver: receiverPubHex,
		Hash:     hash,
	}
	msg, _ := json.Marshal(rm)
	sig, err := utils.CreateSignature(fakePriv, msg)
	assert.Nil(t, err)
	d.Signature = sig
	d.Data = rm
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_SIGNATURE_NOT_VERIFIED, errors.New(resp.Log))
}

func TestReceivedSuccess(t *testing.T) {
	app := NewTFApplication()
	priv, receiverPubHex := utils.GenerateKeyPair()

	hash := createForward(t, app, receiverPubHex)

	d := models.DeliveryModel{}
	d.Action = models.RECEIVED
	rm := models.ReceivedModel{
		Receiver: receiverPubHex,
		Hash:     hash,
	}
	msg, _ := json.Marshal(rm)
	sig, err := utils.CreateSignature(priv, msg)
	assert.Nil(t, err)
	d.Signature = sig
	d.Data = rm
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeOK, resp.Code)

	_, err = app.state.GetForward(hash)
	assert.NotNil(t, err)

	hashes := app.state.GetReceiverHashes(receiverPubHex)
	for _, v := range hashes {
		assert.False(t, v == hash)
	}

	hash = createForward(t, app, receiverPubHex)
	hashes = app.state.GetReceiverHashes(receiverPubHex)
	isFound := false
	for _, v := range hashes {
		if v == hash {
			isFound = true
		}
	}
	assert.True(t, isFound)

}
