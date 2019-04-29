package presale

import (
	"github.com/33cn/plugin/plugin/dapp/presale/executor"
	"github.com/33cn/plugin/plugin/dapp/presale/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.PreSaleX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
