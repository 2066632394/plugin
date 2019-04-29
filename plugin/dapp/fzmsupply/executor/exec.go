package executor

import (
	f "github.com/33cn/plugin/plugin/dapp/fzmsupply/types"
	"github.com/33cn/chain33/types"
)

func (g *Fzmsupply) Exec_InitPlatform(payload *f.RequestInitPlatform, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	uid := g.getTxUid(tx)
	return action.InitPlatform(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_CreateRole(payload *f.RequestCreateRole, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateRole(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_UpdateRole(payload *f.RequestUpdateRole, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.UpdateRole(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_DeleteRole(payload *f.RequestDeleteRole, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.DeleteRole(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_CreateOfficial(payload *f.RequestCreateOfficial, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateOfficial(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_UpdateOfficial(payload *f.RequestUpdateOfficial, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.UpdateOfficial(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_DeleteOfficial(payload *f.RequestDeleteOfficial, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.DeleteOfficial(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_CreateCompany(payload *f.RequestCreateCompany, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateCompany(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_CreateAsset(payload *f.RequestCreateAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateAsset(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_ExamineAsset(payload *f.RequestExamineAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ExamineAsset(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_CancelAsset(payload *f.RequestCancelAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CancelAsset(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_SetCredit(payload *f.RequestSetCredit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.SetCredit(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_ExamineDeposit(payload *f.RequestExamineDeposit, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ExamineDeposit(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_ApplyWithdraw(payload *f.RequestApplyWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ApplyWithdraw(payload,tx.Signature.Pubkey)
}


func (g *Fzmsupply) Exec_ExamineWithdraw(payload *f.RequestExamineWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ExamineWithdraw(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_ReviewWithdraw(payload *f.RequestReviewWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ReviewWithdraw(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Mortgage(payload *f.RequestMortgage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Mortgage(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_CancelMortgage(payload *f.RequestCancelMortgage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CancelMortgage(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_ConfirmMortgage(payload *f.RequestConfirmMortgage, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.ConfirmMortgage(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Redeem(payload *f.RequestRedeem, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateRole(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Overdue(payload *f.RequestOverdue, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateRole(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_ConfirmWithdraw(payload *f.RequestConfirmWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CreateRole(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Receive(payload *f.RequestReceive, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Receive(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Sell(payload *f.RequestSell, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Sell(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Cancel(payload *f.RequestCancel, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Cancel(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_InvestAsset(payload *f.RequestInvestAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.InvestAsset(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Pay(payload *f.RequestPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Pay(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_Recover(payload *f.RequestRecover, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.Recover(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_CancelPay(payload *f.RequestCancelPay, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.CancelPay(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_BatchCash(payload *f.RequestBatchCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.BatchCash(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) Exec_OfflineCash(payload *f.RequestOfflineCash, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := NewAction(g, tx, index)
	return action.OfflineCash(payload,tx.Signature.Pubkey)
}

func (g *Fzmsupply) getTxUid(tx *types.Transaction) string {
	var payload f.Request
	types.Decode(tx.Payload,&payload)
	var uid string
	uid = payload.Uid
	return uid
}