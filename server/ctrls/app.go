package ctrls

import (
	"github.com/mragiadakos/tenderforwarding/server/dbpkg"
	"github.com/tendermint/abci/types"
	dbm "github.com/tendermint/tmlibs/db"
)

var _ types.Application = (*TFApplication)(nil)

type TFApplication struct {
	types.BaseApplication
	state dbpkg.State
}

func NewTFApplication() *TFApplication {
	state := dbpkg.LoadState(dbm.NewMemDB())
	return &TFApplication{state: state}
}
