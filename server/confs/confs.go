package confs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	shell "github.com/ipfs/go-ipfs-api"

	uuid "github.com/satori/go.uuid"
)

type configuration struct {
	AbciDaemon     string
	IpfsConnection string
	Redistributors []string
}

var Conf = configuration{}

func init() {
	Conf.IpfsConnection = "127.0.0.1:5001"
	Conf.AbciDaemon = "tcp://0.0.0.0:26658"
	Conf.Redistributors = []string{}
}

func (c *configuration) SetRedistributorsFromHash(hash string) error {
	sh := shell.NewShell(c.IpfsConnection)
	file := uuid.NewV4().String()
	err := sh.Get(hash, file)
	if err != nil {
		return errors.New("Failed to get the json file for the inflators " + hash + ": " + err.Error())
	}
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.New("Failed to read the json file from the hash " + hash + " for the inflators: " + err.Error())
	}
	os.RemoveAll(file)
	redistributors := []string{}
	err = json.Unmarshal(b, &redistributors)
	if err != nil {
		return errors.New("The json file for the inflators has not the correct JSON format: " + err.Error())
	}
	c.Redistributors = redistributors
	return nil
}
