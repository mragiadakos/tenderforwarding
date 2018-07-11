package ctrls

import (
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/tendermint/abci/types"
)

func (app *TFApplication) DeliverTx(tx []byte) types.ResponseDeliverTx {

	return types.ResponseDeliverTx{Code: models.CodeTypeOK}
}
