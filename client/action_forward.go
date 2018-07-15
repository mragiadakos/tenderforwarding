package main

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/mragiadakos/tenderforwarding/server/models"
	"github.com/mragiadakos/tenderforwarding/server/utils"
	client "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/types"
)

func forward(keyPair KeyPairJson, receiverPubHex, vaultFolder string, coins []string, metadata map[string]string) error {
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

	secret, err := utils.GenerateSharedSecret(priv, receiverPubHex)
	if err != nil {
		return errors.New("Error: could not create a shared secret.")
	}
	encryptedOwners := map[string]string{}
	for _, v := range coins {
		encryptedOwners[v] = ""
	}
	files, err := ioutil.ReadDir(vaultFolder)
	if err != nil {
		return errors.New("Error: " + err.Error())
	}

	for _, v := range files {
		_, ok := encryptedOwners[v.Name()]
		if ok {
			fileB, err := ioutil.ReadFile(vaultFolder + "/" + v.Name())
			if err != nil {
				return errors.New("Error: " + err.Error())
			}
			encryption, err := utils.Encrypt(fileB, secret)
			if err != nil {
				return errors.New("Error: encryption failed: " + err.Error())
			}
			encryptedOwners[v.Name()] = hex.EncodeToString(encryption)
		}
	}
	d := models.DeliveryModel{}
	d.Action = models.FORWARD
	fm := models.ForwardModel{
		Redistributor:   keyPair.PublicKey,
		Coins:           coins,
		EncryptedOwners: encryptedOwners,
		Metadata:        metadata,
		Receiver:        receiverPubHex,
	}
	d.Data = fm
	msg, _ := json.Marshal(d.Data)

	signature, err := utils.CreateSignature(priv, msg)
	d.Signature = signature
	b, _ := json.Marshal(d)

	cli := client.NewHTTP(Confs.TendermintNode, "/websocket")
	btc, err := cli.BroadcastTxCommit(types.Tx(b))
	if err != nil {
		return errors.New("Error: " + err.Error())
	}
	if btc.CheckTx.Code > models.CodeTypeOK {
		return errors.New("Error: " + btc.CheckTx.Log)
	}
	return nil
}
