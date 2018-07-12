package ctrls

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/mragiadakos/tenderforwarding/server/confs"
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
	"github.com/mragiadakos/tenderforwarding/server/validators"

	"github.com/stretchr/testify/assert"
)

func TestForwardFailRedistributorNotInTheList(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor: pubHex,
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_REDISTRIBUTOR_DOES_NOT_EXISTS, errors.New(resp.Log))

}

func TestForwardFailCoinsEmpty(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	confs.Conf.Redistributors = []string{pubHex}

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor: pubHex,
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_COINS_EMPTY, errors.New(resp.Log))

}

func TestForwardFailEncryptedEmpty(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	confs.Conf.Redistributors = []string{pubHex}

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor: pubHex,
		Coins:         []string{uuid.NewV4().String()},
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_ENCRYPTION_EMPTY, errors.New(resp.Log))

}

func TestForwardFailFakeCoin(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	confs.Conf.Redistributors = []string{pubHex}
	coin := uuid.NewV4().String()

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor:   pubHex,
		Coins:           []string{coin},
		EncryptedOwners: map[string]string{uuid.NewV4().String(): "lalalla fake encryption"},
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_COIN_NOT_IN_LIST(coin), errors.New(resp.Log))
}

func TestForwardFailMetadataEmpty(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	confs.Conf.Redistributors = []string{pubHex}
	coin := uuid.NewV4().String()

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor:   pubHex,
		Coins:           []string{coin},
		EncryptedOwners: map[string]string{coin: "lalalla fake encryption"},
	}
	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_METADATA_EMPTY, errors.New(resp.Log))
}

func TestForwardFailNotSignature(t *testing.T) {
	app := NewTFApplication()
	_, pubHex := utils.GenerateKeyPair()
	fakePriv, _ := utils.GenerateKeyPair()

	confs.Conf.Redistributors = []string{pubHex}
	coin := uuid.NewV4().String()

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor:   pubHex,
		Coins:           []string{coin},
		EncryptedOwners: map[string]string{coin: "lalalla fake encryption"},
		Metadata:        map[string]string{"alal": "lalala"},
	}
	msg, _ := json.Marshal(d.Data)
	signature, err := utils.CreateSignature(fakePriv, msg)
	assert.Nil(t, err)
	d.Signature = signature

	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_SIGNATURE_NOT_VERIFIED, errors.New(resp.Log))
}

func TestForwardFailReceiver(t *testing.T) {
	app := NewTFApplication()
	priv, pubHex := utils.GenerateKeyPair()

	confs.Conf.Redistributors = []string{pubHex}
	coin := uuid.NewV4().String()

	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	d.Data = models.ForwardModel{
		Redistributor:   pubHex,
		Coins:           []string{coin},
		EncryptedOwners: map[string]string{coin: "lalalla fake encryption"},
		Metadata:        map[string]string{"alal": "lalala"},
	}
	msg, _ := json.Marshal(d.Data)
	signature, err := utils.CreateSignature(priv, msg)
	assert.Nil(t, err)
	d.Signature = signature

	b, _ := json.Marshal(d)
	resp := app.DeliverTx(b)
	assert.Equal(t, models.CodeTypeUnauthorized, resp.Code)
	assert.Equal(t, validators.ERR_RECEIVER_EMPTY, errors.New(resp.Log))
}

func TestForwardSuccess(t *testing.T) {
	app := NewTFApplication()
	priv, pubHex := utils.GenerateKeyPair()
	_, receiverPubHex := utils.GenerateKeyPair()

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

	nfm, err := app.state.GetForward(hashHex)
	assert.Nil(t, err)

	assert.Equal(t, fm, *nfm)
}
