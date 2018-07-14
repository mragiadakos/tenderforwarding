package ctrls

import (
	"crypto/ecdsa"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
	"github.com/tendermint/abci/types"
)

func receiveForward(t *testing.T, app *TFApplication, priv *ecdsa.PrivateKey, receiverPubHex, hash string) {
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
}

func TestQueryGetForwardsSuccess(t *testing.T) {
	app := NewTFApplication()
	priv, receiverPubHex := utils.GenerateKeyPair()

	createForward(t, app, receiverPubHex)
	hash := createForward(t, app, receiverPubHex)

	qreq := types.RequestQuery{}
	qreq.Path = "get_forwards_by_receiver?pub_hex=" + receiverPubHex

	qresp := app.Query(qreq)
	fwds := []models.ForwardModel{}
	json.Unmarshal(qresp.Value, &fwds)

	assert.Equal(t, 2, len(fwds))
	receiveForward(t, app, priv, receiverPubHex, hash)

	qresp = app.Query(qreq)
	fwds = []models.ForwardModel{}
	json.Unmarshal(qresp.Value, &fwds)
	assert.Equal(t, 1, len(fwds))

}
