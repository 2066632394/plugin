package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	//pty "github.com/33cn/plugin/plugin/dapp/presale/types"
)

var clog = log.New("module", "execs.presale")
var driverName = "presale"

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&PreSale{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register goldbill execer")
	drivers.Register(GetName(), newPreSale, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newPreSale().GetName()
}

type PreSale struct {
	drivers.DriverBase
}

func newPreSale() drivers.Driver {
	n := &PreSale{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *PreSale) GetDriverName() string {
	return driverName
}

func (n *PreSale) CheckTx(tx *types.Transaction, index int) error {
	return nil
}

func Key(str string) (key []byte) {
	key = append(key, []byte("mavl-presale-")...)
	key = append(key, str...)
	return key
}
