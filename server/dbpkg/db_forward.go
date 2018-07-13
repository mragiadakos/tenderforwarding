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
	forwardKey  = []byte("forward:")
	receiverKey = []byte("receiver:")
)

func prefixForward(hashHex string) []byte {
	return append(forwardKey, []byte(hashHex)...)
}

func prefixReceiver(receiver string) []byte {
	return append(receiverKey, []byte(receiver)...)
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

func (s *State) DeleteForward(hashHex string) {
	s.db.Delete(prefixForward(hashHex))
}

func (s *State) AddHashToReceiver(receiver, hash string) {
	b := s.db.Get(prefixReceiver(receiver))
	hashes := []string{}
	json.Unmarshal(b, &hashes)

	hashes = append(hashes, hash)
	bhashes, _ := json.Marshal(hashes)
	s.db.Set(prefixReceiver(receiver), bhashes)
}

func (s *State) RemoveHashFromReceiver(receiver, hash string) {
	b := s.db.Get(prefixReceiver(receiver))
	hashes := []string{}
	json.Unmarshal(b, &hashes)

	newHashes := []string{}
	for _, v := range hashes {
		if v != hash {
			newHashes = append(newHashes, v)
		}
	}

	bhashes, _ := json.Marshal(newHashes)
	s.db.Set(prefixReceiver(receiver), bhashes)
}

func (s *State) GetReceiverHashes(receiver string) []string {
	b := s.db.Get(prefixReceiver(receiver))
	hashes := []string{}
	json.Unmarshal(b, &hashes)
	return hashes
}
