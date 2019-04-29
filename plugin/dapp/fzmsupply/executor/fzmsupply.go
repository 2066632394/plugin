package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	tty "github.com/33cn/plugin/plugin/dapp/fzmsupply/types"
	"github.com/33cn/chain33/common/address"
)

var clog = log.New("module", "execs.fzmsupply")
var driverName = "fzmsupply"
var adminKey = "f82cde5927ce86aab98f2ba388123b56eb165c76a48666d72ade4369f6af18a1"
var conf       = types.ConfSub(driverName)

func init() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Fzmsupply{}))
}

func Init(name string, sub []byte) {
	clog.Debug("register fzmsupply execer")
	drivers.Register(GetName(), newFzmsupply, types.GetDappFork(driverName, "Enable"))
}

func GetName() string {
	return newFzmsupply().GetName()
}

type Fzmsupply struct {
	drivers.DriverBase
}

func newFzmsupply() drivers.Driver {
	n := &Fzmsupply{}
	n.SetChild(n)
	n.SetIsFree(true)
	n.SetExecutorType(types.LoadExecutorType(driverName))
	return n
}

func (n *Fzmsupply) GetDriverName() string {
	return driverName
}

func (n *Fzmsupply) CheckTx(tx *types.Transaction, index int) error {
	var payload tty.Request
	err := types.Decode(tx.Payload,&payload)
	if err != nil {
		return tty.ErrDocodeErr
	}
	if payload.Value == nil {
		return tty.ErrEmptyValue
	}
	switch payload.Ty {
	case tty.FzmsupplyInitPlatformAction:
		return n.checkInitPlatform(payload.GetInitPlatform(),tx.Signature.Pubkey)
	case tty.FzmsupplyCreateRoleAction:
		return n.checkCreateRole(payload.GetCreateRole(),tx.Signature.Pubkey)
	case tty.FzmsupplyUpdateRoleAction:
		return n.checkUpdateRole(payload.GetUpdateRole(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyDeleteRoleAction:
		return n.checkDeleteRole(payload.GetDeleteRole(),tx.Signature.Pubkey)
	case tty.FzmsupplyCreateOfficialAction:
		return n.checkCreateOfficial(payload.GetCreateOfficial(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyUpdateOfficialAction:
		return n.checkUpdateOfficial(payload.GetUpdateOfficial(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyDeleteOfficialAction:
		return n.checkDeleteOfficial(payload.GetDeleteOfficial(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCreateCompanyAction:
		return n.checkCreateCompany(payload.GetCreateCompany(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCreateAssetAction:
		return n.checkCreateAsset(payload.GetCreateAsset(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyExamineAssetAction:
		return n.checkExamineAsset(payload.GetExamineAsset(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCancelAssetAction:
		return n.checkCancelAsset(payload.GetCancelAsset(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplySetCreditAction:
		return n.checkSetCredit(payload.GetSetCredit(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyExamineDepositAction:
		return n.checkExamineDeposit(payload.GetExamineDeposit(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyApplyWithdrawAction:
		return n.checkApplyWithdraw(payload.GetApplyWithdraw(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyExamineWithdrawAction:
		return n.checkExamineWithdraw(payload.GetExamineWithdraw(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyReviewWithdrawAction:
		return n.checkReviewWithdraw(payload.GetReviewWithdraw(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyMortgageAction:
		return n.checkMortgage(payload.GetMortgage(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCancelMortgageAction:
		return n.checkCancelMortgage(payload.GetCancelMortgage(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyConfirmMortgageAction:
		return n.checkConfirmMortgage(payload.GetConfirmMortgage(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyRedeemAction:
		return n.checkRedeem(payload.GetRedeem(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyOverdueAction:
		return n.checkOverdue(payload.GetOverdue(),tx.Signature.Pubkey)
	case tty.FzmsupplyConfirmWithdrawAction:
		return n.checkConfirmWithdraw(payload.GetConfirmWithdraw(),tx.Signature.Pubkey)
	case tty.FzmsupplyReceiveAction:
		return n.checkReceive(payload.GetReceive(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplySellAction:
		return n.checkSell(payload.GetSell(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCancelAction:
		return n.checkCancel(payload.GetCancel(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyInvestAssetAction:
		return n.checkInvestAsset(payload.GetInvestAsset(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyPayAction:
		return n.checkPay(payload.GetPay(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyRecoverAction:
		return n.checkRecover(payload.GetRecover(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyCancelPayAction:
		return n.checkCancelPay(payload.GetCancelPay(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyBatchCashAction:
		return n.checkBatchCash(payload.GetBatchCash(),payload.Uid,tx.Signature.Pubkey)
	case tty.FzmsupplyOfflineCashAction:
		return n.checkOfflineCash(payload.GetOfflineCash (),payload.Uid,tx.Signature.Pubkey)
	default:
		return tty.ErrWrongActionType
	}
	return nil
}


func (n *Fzmsupply) checkAdmin(pubkey []byte) error {

	addr := address.PubKeyToAddress(pubkey).String()
	if ok := IsSuperManager(addr);!ok {
		return tty.ErrWrongPubkey
	}
	return nil
}


// IsSuperManager is supper manager or not
func IsSuperManager(addr string) bool {
	for _, m := range conf.GStrList("superManager") {
		if addr == m {
			return true
		}
	}
	return false
}
