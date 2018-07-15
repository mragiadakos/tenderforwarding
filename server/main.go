package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/mragiadakos/tenderforwarding/server/confs"
	"github.com/mragiadakos/tenderforwarding/server/ctrls"
	absrv "github.com/tendermint/abci/server"
	cmn "github.com/tendermint/tmlibs/common"
	tmlog "github.com/tendermint/tmlibs/log"
)

func main() {
	logger := tmlog.NewTMLogger(kitlog.NewSyncWriter(os.Stdout))
	flagAbci := "socket"

	ipfsDaemon := flag.String("ipfs", "127.0.0.1:5001", "the URL for the IPFS's daemon")
	node := flag.String("node", "tcp://0.0.0.0:26658", "the TCP URL for the ABCI daemon")
	redistributorsHash := flag.String("redistributors", "", "the IPFS hash with the json for the redistributors")
	flag.Parse()

	if len(*redistributorsHash) == 0 {
		fmt.Println("Error:", errors.New("The IPFS hash with the json for the redistributors is missing"))
		return
	}

	err := confs.Conf.SetRedistributorsFromHash(*redistributorsHash)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return
	}
	confs.Conf.AbciDaemon = *node
	confs.Conf.IpfsConnection = *ipfsDaemon

	app := ctrls.NewTFApplication()
	srv, err := absrv.NewServer(confs.Conf.AbciDaemon, flagAbci, app)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		srv.Stop()
	})

}

/*
func startTendermint() error {
	tcmd.AddNodeFlags(tcmd.RootCmd)
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}
	log.Fatal(cfg)

	logger := tmlog.NewTMLogger(kitlog.NewSyncWriter(os.Stdout))

	app := ctrls.NewTFApplication()

	n, err := node.NewNode(cfg,
		privval.LoadOrGenFilePV(cfg.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		logger,
	)
	if err != nil {
		return fmt.Errorf("Failed to create node: %v", err)
	}

	if err := n.Start(); err != nil {
		return fmt.Errorf("Failed to start node: %v", err)
	}
	logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
*/
