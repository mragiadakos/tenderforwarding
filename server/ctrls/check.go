package ctrls

import (
	"encoding/json"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/validators"
	"github.com/tendermint/abci/types"
)

func (app *TFApplication) CheckTx(tx []byte) types.ResponseCheckTx {
	dm := models.DeliveryModel{}
	err := json.Unmarshal(tx, &dm)
	if err != nil {
		return types.ResponseCheckTx{Code: models.CodeTypeEncodingError, Log: "The delivery is not correct json: " + err.Error()}
	}
	switch dm.Action {
	case models.FORWARD:
		fm := dm.GetForwardModel()
		status, err := validators.ValidateForward(&fm, dm.Signature)
		if err != nil {
			return types.ResponseCheckTx{Code: status, Log: err.Error()}
		}
	case models.RECEIVED:
		rm := dm.GetReceivedModel()
		status, err := validators.ValidateReceived(&app.state, &rm, dm.Signature)
		if err != nil {
			return types.ResponseCheckTx{Code: status, Log: err.Error()}
		}

	default:
		return types.ResponseCheckTx{Code: models.CodeTypeUnauthorized, Log: ERR_ACTION_DOES_NOT_EXIST.Error()}
	}

	return types.ResponseCheckTx{Code: models.CodeTypeOK}
}
