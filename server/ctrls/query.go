package ctrls

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/tendermint/abci/types"
)

var (
	ERR_THE_QUERY_METHOD_HAS_NOT_BEEN_FOUND = errors.New("The query's method has not been found.")
	ERR_PUBLIC_KEY_HAS_NOT_BEEN_SUBMITTED   = errors.New("The public key has not been submitted.")
)

func getForwardsByReceiver(app *TFApplication, pubHex string) []models.ForwardModel {
	hashes := app.state.GetReceiverHashes(pubHex)
	fwds := []models.ForwardModel{}
	for _, hash := range hashes {
		fwd, err := app.state.GetForward(hash)
		if err == nil {
			fwds = append(fwds, *fwd)
		}
	}
	return fwds
}

func (app *TFApplication) Query(qreq types.RequestQuery) types.ResponseQuery {
	u, err := url.Parse(qreq.Path)
	if err != nil {
		return types.ResponseQuery{Code: models.CodeTypeEncodingError, Log: err.Error()}
	}

	switch u.Path {
	case "get_forwards_by_receiver":
		values := u.Query()
		pubHex := values.Get("pub_hex")
		if len(pubHex) == 0 {
			return types.ResponseQuery{Code: models.CodeTypeUnauthorized,
				Log: ERR_PUBLIC_KEY_HAS_NOT_BEEN_SUBMITTED.Error()}
		}
		fwds := getForwardsByReceiver(app, pubHex)
		b, _ := json.Marshal(fwds)
		return types.ResponseQuery{Code: models.CodeTypeOK, Value: b}
	}

	return types.ResponseQuery{Code: models.CodeTypeUnauthorized,
		Log: ERR_THE_QUERY_METHOD_HAS_NOT_BEEN_FOUND.Error()}

}
