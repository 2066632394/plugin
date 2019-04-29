package executor

import (
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

func (c *Contract) Exec_Create(payload *cty.ContractCreate, tx *types.Transaction, index int) (*types.Receipt, error) {
	return newContractAction(c, tx).contractCreate(payload)
}

func (c *Contract) Exec_Cancel(payload *cty.ContractCancel, tx *types.Transaction, index int) (*types.Receipt, error) {
	return newContractAction(c, tx).contractCancel(payload)
}

func (c *Contract) Exec_Modify(payload *cty.ContractModify, tx *types.Transaction, index int) (*types.Receipt, error) {
	return newContractAction(c, tx).contractModify(payload)
}

func (c *Contract) Exec_Sign(payload *cty.ContractSign, tx *types.Transaction, index int) (*types.Receipt, error) {
	return newContractAction(c, tx).contractSign(payload)
}

func (c *Contract) Exec_Reject(payload *cty.ContractReject, tx *types.Transaction, index int) (*types.Receipt, error) {
	return newContractAction(c, tx).contractReject(payload)
}

//func (c *Contract) Exec_Register(payload *cty.ContractRegister, tx *types.Transaction, index int) (*types.Receipt, error) {
//	return newContractAction(c, tx).userRegister(payload)
//}
