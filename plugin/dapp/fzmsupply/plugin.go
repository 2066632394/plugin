package fzmsupply

import (
	"github.com/33cn/plugin/plugin/dapp/fzmsupply/executor"
	"github.com/33cn/plugin/plugin/dapp/fzmsupply/types"
	"github.com/33cn/chain33/pluginmgr"
)

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     types.FzmsupplyX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      nil,
		RPC:      nil,
	})
}
