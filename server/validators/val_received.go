package validators

import (
	"encoding/json"
	"errors"

	"github.com/mragiadakos/tenderforwarding/server/dbpkg"
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
)

var (
	ERR_HASH_DOES_NOT_EXISTS    = errors.New("The hash does not exists.")
	ERR_RECEIVER_NOT_IN_FORWARD = errors.New("The receiver is not in the forward.")
)

func ValidateReceived(state *dbpkg.State, rm *models.ReceivedModel, sig string) (uint32, error) {
	fm, err := state.GetForward(rm.Hash)
	if err != nil {
		return models.CodeTypeUnauthorized, ERR_HASH_DOES_NOT_EXISTS
	}

	if fm.Receiver != rm.Receiver {
		return models.CodeTypeUnauthorized, ERR_RECEIVER_NOT_IN_FORWARD
	}
	msg, _ := json.Marshal(rm)
	ver, err := utils.Verify(rm.Receiver, sig, msg)
	if err != nil {
		return models.CodeTypeEncodingError, err
	}

	if !ver {
		return models.CodeTypeUnauthorized, ERR_SIGNATURE_NOT_VERIFIED
	}
	return models.CodeTypeOK, nil
}
