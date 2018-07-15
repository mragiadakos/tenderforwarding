package main

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mragiadakos/tenderforwarding/server/utils"
	"github.com/urfave/cli"
)

type KeyPairJson struct {
	PrivateKey string
	PublicKey  string
}

var GenerateKeyCommand = cli.Command{
	Name:    "generate_key",
	Aliases: []string{"g"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "filename",
			Usage: "the filename that the key will be saved",
		},
	},
	Usage: "generate the key in a file",
	Action: func(c *cli.Context) error {
		filename := c.String("filename")
		if len(filename) == 0 {
			return errors.New("Error: filename is missing")
		}

		priv, pubHex := utils.GenerateKeyPair()
		b, err := x509.MarshalPKCS8PrivateKey(priv)
		if err != nil {
			return errors.New("Error: Could not marshal private key: " + err.Error())
		}

		privHex := hex.EncodeToString(b)

		kp := KeyPairJson{
			PublicKey:  pubHex,
			PrivateKey: privHex,
		}

		kpb, _ := json.Marshal(kp)
		err = ioutil.WriteFile(filename, kpb, 0644)
		if err != nil {
			return errors.New("Error: failed to write the file: " + err.Error())
		}
		fmt.Println("The generate was successful")
		return nil
	},
}

var ForwardCommand = cli.Command{
	Name:    "forward",
	Aliases: []string{"f"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
		cli.StringFlag{
			Name:  "vault",
			Usage: "the folder of the vault with the coins",
		},
		cli.StringFlag{
			Name:  "receiver",
			Usage: "the receiver's public key",
		},
		cli.StringFlag{
			Name:  "coins",
			Usage: "the list of coins seperated by comma",
		},
		cli.StringFlag{
			Name:  "reason",
			Usage: "the reason of the forward",
		},
	},
	Usage: "Forward coins",
	Action: func(c *cli.Context) error {
		keyFilename := c.String("key")
		if len(keyFilename) == 0 {
			return errors.New("Error: key is missing")
		}

		vaultFolder := c.String("vault")
		if len(vaultFolder) == 0 {
			return errors.New("Error: vault is missing")
		}

		coinsString := c.String("coins")
		if len(coinsString) == 0 {
			return errors.New("Error: coins are missing")
		}

		reason := c.String("reason")
		if len(reason) == 0 {
			return errors.New("Error: reason is missing")
		}

		receiver := c.String("receiver")
		if len(receiver) == 0 {
			return errors.New("Error: the receiver is missing")
		}

		keyPairBytes, err := ioutil.ReadFile(keyFilename)
		if err != nil {
			return errors.New("Error: key's filename failed: " + err.Error())
		}
		keyPair := KeyPairJson{}
		err = json.Unmarshal(keyPairBytes, &keyPair)
		if err != nil {
			return errors.New("Error: key's filename does not decode: " + err.Error())
		}

		err = forward(keyPair, receiver, vaultFolder, strings.Split(coinsString, ","), map[string]string{"reason": reason})
		if err != nil {
			return err
		}

		fmt.Println("The forward was successfull.")
		return nil
	},
}

var QueryCommand = cli.Command{
	Name:    "query",
	Aliases: []string{"q"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
	},
	Usage: "query forwards",
	Action: func(c *cli.Context) error {
		keyFilename := c.String("key")
		if len(keyFilename) == 0 {
			return errors.New("Error: key is missing")
		}

		keyPairBytes, err := ioutil.ReadFile(keyFilename)
		if err != nil {
			return errors.New("Error: key's filename failed: " + err.Error())
		}
		keyPair := KeyPairJson{}
		err = json.Unmarshal(keyPairBytes, &keyPair)
		if err != nil {
			return errors.New("Error: key's filename does not decode: " + err.Error())
		}

		err = query(keyPair)
		if err != nil {
			return err
		}

		return nil
	},
}

var ReceiveCommand = cli.Command{
	Name:    "receive",
	Aliases: []string{"r"},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "key",
			Usage: "the filename of the key",
		},
		cli.StringFlag{
			Name:  "vault",
			Usage: "the folder of the vault that we will add the coins",
		},
		cli.StringFlag{
			Name:  "hash",
			Usage: "the hash of the forward with the coins",
		},
	},
	Usage: "receive forwards",
	Action: func(c *cli.Context) error {
		keyFilename := c.String("key")
		if len(keyFilename) == 0 {
			return errors.New("Error: key is missing")
		}

		vaultFolder := c.String("vault")
		if len(vaultFolder) == 0 {
			return errors.New("Error: vault is missing")
		}

		hash := c.String("hash")
		if len(hash) == 0 {
			return errors.New("Error: hash is missing")
		}

		keyPairBytes, err := ioutil.ReadFile(keyFilename)
		if err != nil {
			return errors.New("Error: key's filename failed: " + err.Error())
		}
		keyPair := KeyPairJson{}
		err = json.Unmarshal(keyPairBytes, &keyPair)
		if err != nil {
			return errors.New("Error: key's filename does not decode: " + err.Error())
		}

		err = receive(keyPair, vaultFolder, hash)
		if err != nil {
			return err
		}

		fmt.Println("The receive was successful.")
		return nil
	},
}
