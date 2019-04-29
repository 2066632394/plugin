package contract

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/contract/executor"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     cty.ContractX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
