package validators

import (
	"encoding/json"
	"errors"

	"github.com/mragiadakos/tenderforwarding/server/confs"
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
)

var (
	ERR_REDISTRIBUTOR_DOES_NOT_EXISTS = errors.New("The redistributor does not exists.")
	ERR_COINS_EMPTY                   = errors.New("The coins are empty.")
	ERR_ENCRYPTION_EMPTY              = errors.New("The encrypted private keys for owners are empty.")
	ERR_COIN_NOT_IN_LIST              = func(coin string) error {
		return errors.New("The coin " + coin + " does not exist in the list of owners.")
	}
	ERR_METADATA_EMPTY         = errors.New("The metadata is empty.")
	ERR_RECEIVER_EMPTY         = errors.New("The receiver is empty.")
	ERR_SIGNATURE_NOT_VERIFIED = errors.New("The signature does not verify the data.")
)

func ValidateForward(fm *models.ForwardModel, sig string) (uint32, error) {
	redistributorExists := false
	for _, v := range confs.Conf.Redistributors {
		if v == fm.Redistributor {
			redistributorExists = true
		}
	}
	if !redistributorExists {
		return models.CodeTypeUnauthorized, ERR_REDISTRIBUTOR_DOES_NOT_EXISTS
	}
	if len(fm.Coins) == 0 {
		return models.CodeTypeUnauthorized, ERR_COINS_EMPTY
	}
	if len(fm.EncryptedOwners) == 0 {
		return models.CodeTypeUnauthorized, ERR_ENCRYPTION_EMPTY
	}

	for _, coin := range fm.Coins {
		_, ok := fm.EncryptedOwners[coin]
		if !ok {
			return models.CodeTypeUnauthorized, ERR_COIN_NOT_IN_LIST(coin)
		}
	}

	if len(fm.Metadata) == 0 {
		return models.CodeTypeUnauthorized, ERR_METADATA_EMPTY
	}

	msg, _ := json.Marshal(fm)
	verifies, err := utils.Verify(fm.Redistributor, sig, msg)
	if err != nil {
		return models.CodeTypeEncodingError, err
	}
	if !verifies {
		return models.CodeTypeUnauthorized, ERR_SIGNATURE_NOT_VERIFIED
	}
	if len(fm.Receiver) == 0 {
		return models.CodeTypeUnauthorized, ERR_RECEIVER_EMPTY
	}

	return models.CodeTypeOK, nil
}
