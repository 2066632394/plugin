package executor

import (
	pty "github.com/33cn/plugin/plugin/dapp/presale/types"
	"github.com/33cn/chain33/types"
)

func (p *PreSale) Exec_MarketReg(payload *pty.PreSaleMarketReg, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execMarketReg(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_AdminReg(payload *pty.PreSaleAdminReg, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execAdminReg(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserReg(payload *pty.PreSaleUserReg, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserReg(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_MarketRec(payload *pty.PreSaleMarketRec, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execMarketRec(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserRec(payload *pty.PreSaleUserRec, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserRec(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_AddPdInfo(payload *pty.PreSaleAddPdInfo, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execAddPdInfo(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_OnOffPd(payload *pty.PreSaleOnOffPd, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execGetOnOffPd(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserPay(payload *pty.PreSaleUserPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserPay(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_Address(payload *pty.PreSaleAddAddress, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execAddAddress(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_OrderTrans(payload *pty.PreSaleOrderTrans, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execOrderTrans(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_MarketShip(payload *pty.PreSaleMarketShip, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execMarketShip(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_Logistics(payload *pty.PreSaleLogistics, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execLogistics(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserRecipt(payload *pty.PreSaleUserRecipt, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserRecipt(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_OverDue(payload *pty.PreSaleOverDue, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execOverDue(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_WithDraw(payload *pty.PreSaleWithDraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execWithDraw(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_ChangePd(payload *pty.PreSaleChangePd, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execChange(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserTokenWithDraw(payload *pty.PreSaleUserTokenWithDraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserTokenWithDraw(payload,tx.Signature.Pubkey)
}

func (p *PreSale) Exec_UserShip(payload *pty.PreSaleUserShip, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(p, tx, index)
	return action.execUserShip(payload,tx.Signature.Pubkey)
}

