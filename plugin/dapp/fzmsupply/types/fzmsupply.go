package types

import (
	"github.com/33cn/chain33/types"

)

var FzmsupplyX = "fzmsupply"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(FzmsupplyX))
	types.RegistorExecutor(FzmsupplyX, NewType())
	types.RegisterDappFork(FzmsupplyX, "Enable", 0)
}

type FzmsupplyType struct {
	types.ExecTypeBase
}

func NewType() *FzmsupplyType {
	c := &FzmsupplyType{}
	c.SetChild(c)
	return c
}

func (g *FzmsupplyType) GetPayload() types.Message {
	return &Request{}
}

func (g *FzmsupplyType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"InitPlatform":FzmsupplyInitPlatformAction,
		"CreateRole":FzmsupplyCreateRoleAction,
		"UpdateRole":FzmsupplyUpdateRoleAction,
		"DeleteRole":FzmsupplyDeleteRoleAction,
		"CreateOfficial":FzmsupplyCreateOfficialAction,
		"UpdateOfficial":FzmsupplyUpdateOfficialAction,
		"DeleteOfficial":FzmsupplyDeleteOfficialAction,
		"CreateCompany":FzmsupplyCreateCompanyAction,
		"CreateAsset":FzmsupplyCreateAssetAction,
		"ExamineAsset":FzmsupplyExamineAssetAction,
		"CancelAsset":FzmsupplyCancelAssetAction,
		"SetCredit":FzmsupplySetCreditAction,
		"ExamineDeposit":FzmsupplyExamineDepositAction,
		"ApplyWithdraw":FzmsupplyApplyWithdrawAction,
		"ExamineWithdraw":FzmsupplyExamineWithdrawAction,
		"ReviewWithdraw":FzmsupplyReviewWithdrawAction,
		"Mortgage":FzmsupplyMortgageAction,
		"CancelMortgage":FzmsupplyCancelMortgageAction,
		"ConfirmMortgage":FzmsupplyConfirmMortgageAction,
		"Redeem":FzmsupplyRedeemAction,
		"Overdue":FzmsupplyOverdueAction,
		"ConfirmWithdraw":FzmsupplyConfirmWithdrawAction,
		"Receive":FzmsupplyReceiveAction,
		"Sell":FzmsupplySellAction,
		"Cancel":FzmsupplyCancelAction,
		"InvestAsset":FzmsupplyInvestAssetAction,
		"Pay":FzmsupplyPayAction,
		"Recover":FzmsupplyRecoverAction,
		"CancelPay":FzmsupplyCancelPayAction,
		"BatchCash":FzmsupplyBatchCashAction,
		"OfflineCash":FzmsupplyOfflineCashAction,
	}
}

func (g *FzmsupplyType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{}
}
