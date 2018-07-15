package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mragiadakos/tenderforwarding/server/dbpkg"
	"github.com/mragiadakos/tenderforwarding/server/models"

	client "github.com/tendermint/tendermint/rpc/client"
)

func query(keyPair KeyPairJson) error {
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

	for _, v := range stfs {
		fmt.Println("Hash:\t", v.Hash)
		fmt.Println("Coins:\t", v.Coins)
		fmt.Println("\n")
	}

	return nil
}
