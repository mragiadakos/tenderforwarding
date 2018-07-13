package ctrls

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/validators"

	"github.com/tendermint/abci/types"
)

var (
	ERR_ACTION_DOES_NOT_EXIST = errors.New("The action does not exists.")
)

func (app *TFApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {
	dm := models.DeliveryModel{}
	err := json.Unmarshal(tx, &dm)
	if err != nil {
		return types.ResponseDeliverTx{Code: models.CodeTypeEncodingError, Log: "The delivery is not correct json: " + err.Error()}
	}
	switch dm.Action {
	case models.FORWARD:
		fm := dm.GetForwardModel()
		status, err := validators.ValidateForward(&fm, dm.Signature)
		if err != nil {
			return types.ResponseDeliverTx{Code: status, Log: err.Error()}
		}
		app.state.AddForward(fm)
		cb, _ := json.Marshal(fm.Coins)
		hash := sha256.Sum256(cb)
		hashHex := hex.EncodeToString(hash[:])
		app.state.AddHashToReceiver(fm.Receiver, hashHex)
	case models.RECEIVED:
		rm := dm.GetReceivedModel()
		status, err := validators.ValidateReceived(&app.state, &rm, dm.Signature)
		if err != nil {
			return types.ResponseDeliverTx{Code: status, Log: err.Error()}
		}
		app.state.DeleteForward(rm.Hash)
		app.state.RemoveHashFromReceiver(rm.Receiver, rm.Hash)
	default:
		return types.ResponseDeliverTx{Code: models.CodeTypeUnauthorized, Log: ERR_ACTION_DOES_NOT_EXIST.Error()}
	}
	return types.ResponseDeliverTx{Code: models.CodeTypeOK}
}
