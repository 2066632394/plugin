package executor

import (
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

func (c *Contract) Query_GetContractInfo(in *cty.ReqId) (types.Message, error) {
	contractLog.Debug("Contract Query", "contractId", in.GetId())

	value, err := c.GetStateDB().Get(calcContractKey(in.GetId()))
	if err != nil {
		return nil, err
	}
	var info cty.Contract
	err = types.Decode(value, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

//func (c *Contract) Query_GetUserInfo(in *cty.ReqId) (types.Message, error) {
//	contractLog.Debug("Contract Query", "contractId", in.GetId())
//
//	value, err := c.GetStateDB().Get(calcUserKey(in.GetId()))
//	if err != nil {
//		return nil, err
//	}
//	var info cty.User
//	err = types.Decode(value, &info)
//	if err != nil {
//		return nil, err
//	}
//	return &info, nil
//}
