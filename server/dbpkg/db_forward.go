package dbpkg

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/mragiadakos/tenderforwarding/server/models"
)

var (
	ERR_FORWARD_DOES_NOT_EXIST = errors.New("The forward does not exists.")
)

var (
	forwardKey = []byte("forward:")
)

func prefixForward(hashHex string) []byte {
	return append(forwardKey, []byte(hashHex)...)
}

func (s *State) AddForward(fm models.ForwardModel) {
	b, _ := json.Marshal(fm.Coins)
	hash := sha256.Sum256(b)
	hashHex := hex.EncodeToString(hash[:])
	fmb, _ := json.Marshal(fm)
	s.db.Set(prefixForward(hashHex), fmb)
}

func (s *State) GetForward(hashHex string) (*models.ForwardModel, error) {
	has := s.db.Has(prefixForward(hashHex))
	if !has {
		return nil, ERR_FORWARD_DOES_NOT_EXIST
	}

	fmb := s.db.Get(prefixForward(hashHex))

	fm := models.ForwardModel{}
	json.Unmarshal(fmb, &fm)
	return &fm, nil
}
