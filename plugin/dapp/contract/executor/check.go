package executor

import (
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
	"strings"
)

type ContractActionValue interface {
	GetOperatorId() string
}

func (c *Contract) CheckTx(tx *types.Transaction, index int) error {
	var action cty.ContractAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		return err
	}
	pubKey := common.Bytes2Hex(tx.GetSignature().GetPubkey())
	switch {
	case action.Ty == cty.ContractActionCreate && action.GetCreate() != nil:
		return c.checkUserPubKey(pubKey, action.GetCreate())
	case action.Ty == cty.ContractActionCancel && action.GetCancel() != nil:
		return c.checkUserPubKey(pubKey, action.GetCancel())
	case action.Ty == cty.ContractActionModify && action.GetModify() != nil:
		return c.checkUserPubKey(pubKey, action.GetModify())
	case action.Ty == cty.ContractActionSign && action.GetSign() != nil:
		return c.checkUserPubKey(pubKey, action.GetSign())
	case action.Ty == cty.ContractActionReject && action.GetReject() != nil:
		return c.checkUserPubKey(pubKey, action.GetReject())
		//case action.Ty == cty.ContractActionRegister && action.GetRegister() != nil:
		//	return c.checkManagerPubKey(pubKey)
	}
	return types.ErrActionNotSupport
}

//func (c *Contract) checkManagerPubKey(pubKey string) error {
//	if strings.EqualFold(pubKey, cty.ManagerPubKey) {
//		return nil
//	}
//	return cty.ErrPubKey
//}

func (c *Contract) checkUserPubKey(pubKey string, actionValue ContractActionValue) error {
	value, err := c.GetStateDB().Get(calcUserKey(actionValue.GetOperatorId()))
	if err != nil {
		return err
	}
	if !strings.EqualFold(pubKey, common.Bytes2Hex(value)) {
		return cty.ErrPubKey
	}
	return nil
}
