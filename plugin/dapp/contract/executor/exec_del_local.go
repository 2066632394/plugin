package executor

import (
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

func (c *Contract) ExecDelLocal_Create(payload *cty.ContractCreate, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{KV: nil}, nil
}

func (c *Contract) ExecDelLocal_Cancel(payload *cty.ContractCancel, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{KV: nil}, nil
}

func (c *Contract) ExecDelLocal_Modify(payload *cty.ContractModify, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{KV: nil}, nil
}

func (c *Contract) ExecDelLocal_Sign(payload *cty.ContractSign, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{KV: nil}, nil
}

func (c *Contract) ExecDelLocal_Reject(payload *cty.ContractReject, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	return &types.LocalDBSet{KV: nil}, nil
}

//func (c *Contract) ExecDelLocal_Register(payload *cty.ContractRegister, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {
//	return &types.LocalDBSet{KV: nil}, nil
//}
