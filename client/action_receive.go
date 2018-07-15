package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/mragiadakos/tenderforwarding/server/dbpkg"
	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
	client "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

func receive(keyPair KeyPairJson, vault, hash string) error {
	cli := client.NewHTTP(Confs.TendermintNode, "/websocket")
	q, err := cli.ABCIQuery("get_forwards_by_receiver?pub_hex="+keyPair.PublicKey, nil)
	if err != nil {
		return errors.New("Error:" + err.Error())
	}
	if q.Response.Code > models.CodeTypeOK {
		return errors.New("Error: " + q.Response.Log)
	}
	stfs := []dbpkg.ForwardState{}
	json.Unmarshal(q.Response.Value, &stfs)

	var st *dbpkg.ForwardState
	for _, v := range stfs {
		if v.Hash == hash {
			st = &v
			break
		}
	}

	if st == nil {
		return errors.New("Error: the hash has not been found.")
	}

	privCert, err := hex.DecodeString(keyPair.PrivateKey)
	if err != nil {
		return errors.New("Error: private key failed on decoding: " + err.Error())
	}

	privInterface, err := x509.ParsePKCS8PrivateKey(privCert)
	if err != nil {
		return errors.New("Error: private key failed on parsing: " + err.Error())
	}

	priv, ok := privInterface.(*ecdsa.PrivateKey)
	if !ok {
		return errors.New("Error: private key is not ecdsa.PrivateKey")
	}

	secret, err := utils.GenerateSharedSecret(priv, st.Redistributor)
	if err != nil {
		return errors.New("Error: could not create a shared secret.")
	}

	for k, v := range st.EncryptedOwners {
		b, err := hex.DecodeString(v)
		if err != nil {
			return errors.New("Error: the coin " + k + "'s the private key has not encoded correctly in hex: " + err.Error())
		}
		decrypted, err := utils.Decrypt(b, secret)
		if err != nil {
			return errors.New("Error: the coin " + k + " has not been decrypted: " + err.Error())
		}

		err = ioutil.WriteFile(vault+"/"+k, decrypted, 0644)
		if err != nil {
			return errors.New("Error: failed to write: " + err.Error())
		}
	}

	d := models.DeliveryModel{}
	d.Action = models.RECEIVED
	rm := models.ReceivedModel{
		Receiver: keyPair.PublicKey,
		Hash:     hash,
	}
	msg, _ := json.Marshal(rm)
	sig, err := utils.CreateSignature(priv, msg)
	d.Signature = sig
	d.Data = rm
	b, _ := json.Marshal(d)

	btc, err := cli.BroadcastTxCommit(types.Tx(b))
	if err != nil {
		return errors.New("Error: " + err.Error())
	}
	if btc.CheckTx.Code > models.CodeTypeOK {
		return errors.New("Error: " + btc.CheckTx.Log)
	}

	return nil
}
