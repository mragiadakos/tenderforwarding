package ctrls

import (
	"encoding/json"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/validators"

	"github.com/tendermint/abci/types"
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
	}
	return types.ResponseDeliverTx{Code: models.CodeTypeOK}
}
