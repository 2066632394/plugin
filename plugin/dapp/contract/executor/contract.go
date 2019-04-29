package executor

import (
	log "github.com/inconshreveable/log15"
	_ "github.com/33cn/chain33/common/db"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

var contractLog = log.New("module", "execs.contract")
var driverName = cty.ContractX

func Init(name string, sub []byte) {
	driverName = name
	cty.ContractX = driverName
	cty.ExecerContract = []byte(driverName)
	drivers.Register(driverName, newContract, 0)
}

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Contract{}))
}

func GetName() string {
	return newContract().GetName()
}

type Contract struct {
	drivers.DriverBase
}

func newContract() drivers.Driver {
	n := &Contract{}
	n.SetChild(n)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	n.SetIsFree(true)
	return n
}

func (c *Contract) GetDriverName() string {
	return driverName
}
