package executor

import (
	g "github.com/33cn/plugin/plugin/dapp/fzmsupply/types"
	"github.com/33cn/chain33/types"
	"bytes"
	"github.com/33cn/chain33/common/address"
)

func (app *Fzmsupply) checkPlatform (pubkey []byte) error {
	plat,err := app.getPlatform()
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, plat.Pubkey) {
		return g.ErrAuthenticate
	}
	return nil
}

func (app *Fzmsupply) checkOperateOfficial(pubkey []byte, uid string) (error) {
	if uid != "" {
		admin,err := app.getOfficial(uid)
		if err != nil {
			return err
		}
		if admin.IsAdmin == false || admin.Active == false {
			return g.ErrNoRight
		}
		if !bytes.Equal(pubkey, admin.Pubkey) {
			return g.ErrAuthenticate
		}

		err = app.checkPlatform(pubkey)
		if err != nil {
			return err
		}
		return nil
	}

	return g.ErrEmptyValue
}

func (app *Fzmsupply) checkIdentity(pubkey []byte, uid string) error {
	user,err := app.getUser(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, user.Pubkey) {
		return g.ErrAuthenticate
	}
	return nil
}


func (app *Fzmsupply) checkInitPlatform (req *g.RequestInitPlatform,pubkey []byte) error {
	addr := address.PubKeyToAddress(pubkey).String()
	if !IsSuperManager(addr) {
		return g.ErrNoRight
	}

	plat,err := app.getPlatform()
	if err == nil && plat != nil {
		return g.ErrPlatformExist
	}
	return nil
}

func (app *Fzmsupply) checkCreateRole (req *g.RequestCreateRole,pubkey []byte) error {

	return app.checkPlatform(pubkey)
}

func (app *Fzmsupply) checkDeleteRole (req *g.RequestDeleteRole,pubkey []byte) error {

	if len(req.Roleids) == 0 {
		return g.ErrNoRecord
	}

	return app.checkPlatform(pubkey)
}

func (app *Fzmsupply) checkOverdue (req *g.RequestOverdue,pubkey []byte) error {

	if len(req.MorgageIds) == 0 {
		return g.ErrNoRecord
	}
	return app.checkPlatform(pubkey)
}

func (app *Fzmsupply) checkConfirmWithdraw (req *g.RequestConfirmWithdraw,pubkey []byte) error {

	if len(req.Withdraws) == 0 {
		return g.ErrNoRecord
	}
	return app.checkPlatform(pubkey)
}


func (app *Fzmsupply) checkCreateOfficial (req *g.RequestCreateOfficial,uid string,pubkey []byte) error {

	return app.checkOperateOfficial(pubkey, uid)
}

func (app *Fzmsupply) checkUpdateOfficial (req *g.RequestUpdateOfficial,uid string,pubkey []byte) error {

	return app.checkOperateOfficial(pubkey, uid)
}

func (app *Fzmsupply) checkDeleteOfficial (req *g.RequestDeleteOfficial,uid string,pubkey []byte) error {

	if len(req.Uids) == 0 {
		return g.ErrNoRecord
	}

	return app.checkOperateOfficial(pubkey, uid)
}

func (app *Fzmsupply) checkUpdateRole (req *g.RequestUpdateRole,uid string,pubkey []byte) error {

	if req.Roleid == "admin" {
		return g.ErrNoRight
	}

	return app.checkOperateOfficial(pubkey, uid)
}

func (app *Fzmsupply) checkCreateCompany(req *g.RequestCreateCompany,uid string,pubkey []byte) error {

	_,err := app.getPlatform()
	if err != nil && err == types.ErrNotFound {
		return g.ErrPlatformNotExist
	}
	user,err := app.getUser(uid)
	if err == nil && user != nil {
		return g.ErrUserExist
	}
	return nil
}

func (app *Fzmsupply) checkCreateAsset(req *g.RequestCreateAsset,uid string,pubkey []byte) error {

	if req.Amount <= 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkExamineAsset(req *g.RequestExamineAsset,uid string,pubkey []byte) error {

	if req.ValidCredit < 0 {
		return g.ErrIllegalAmount
	}

	err := app.checkPlatform(pubkey)
	if err == nil {
		return nil
	}
	official,err := app.getOfficial(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, official.Pubkey) {
		return g.ErrAuthenticate
	}
	if official.Active == false {
		return g.ErrNoRight
	}

	var right int32 = 0
	for _, v := range official.Roleids {
		role ,err := app.getRole(v)
		if err != nil {
			return err
		}
		right |= role.Rights
	}

	if right & int32(g.RightType_RsExamineAsset) == 0 {
		return g.ErrNoRight
	}
	return nil
}

func (app *Fzmsupply) checkCancelAsset(req *g.RequestCancelAsset,uid string,pubkey []byte) error {

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkSetCredit (req *g.RequestSetCredit,uid string,pubkey []byte) error {
	err := app.checkPlatform(pubkey)
	if err == nil {
		return nil
	}
	official,err := app.getOfficial(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, official.Pubkey) {
		return g.ErrAuthenticate
	}
	if official.Active == false {
		return g.ErrNoRight
	}
	var right int32 = 0
	for _, v := range official.Roleids {
		role ,err := app.getRole(v)
		if err != nil {
			return err
		}
		right |= role.Rights
	}

	if right & int32(g.RightType_RsSetCredit) == 0 {
		return g.ErrNoRight
	}

	return nil
}

func (app *Fzmsupply) checkExamineDeposit (req *g.RequestExamineDeposit,uid string,pubkey []byte) error {

	if len(req.Deposits) == 0 {
		return g.ErrNoRecord
	}

	err := app.checkPlatform(pubkey)
	if err == nil {
		return nil
	}
	official,err := app.getOfficial(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, official.Pubkey) {
		return g.ErrAuthenticate
	}
	if official.Active == false {
		return g.ErrNoRight
	}

	var right int32 = 0
	for _, v := range official.Roleids {
		role ,err := app.getRole(v)
		if err != nil {
			return err
		}
		right |= role.Rights
	}

	if right & int32(g.RightType_RsExamineDeposit) == 0 {
		return g.ErrNoRight
	}

	return nil
}

func (app *Fzmsupply) checkApplyWithdraw (req *g.RequestApplyWithdraw,uid string,pubkey []byte) error {

	with ,err := app.getWithdraw(req.GetWithdrawId())
	if err == nil && with != nil  {
		return g.ErrWithdrawExist
	}
	if req.Rmb <= 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkExamineWithdraw (req *g.RequestExamineWithdraw,uid string,pubkey []byte) error {

	if len(req.Withdraws) == 0 {
		return g.ErrNoRecord
	}

	err := app.checkPlatform(pubkey)
	if err == nil {
		return nil
	}
	official,err := app.getOfficial(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, official.Pubkey) {
		return g.ErrAuthenticate
	}
	if official.Active == false {
		return g.ErrNoRight
	}

	var right int32 = 0
	for _, v := range official.Roleids {
		role ,err := app.getRole(v)
		if err != nil {
			return err
		}
		right |= role.Rights
	}

	if right & int32(g.RightType_RsExamineWithdraw) == 0 {
		return g.ErrNoRight
	}

	return nil
}

func (app *Fzmsupply) checkReviewWithdraw (req *g.RequestReviewWithdraw,uid string,pubkey []byte) error {

	if len(req.Withdraws) == 0 {
		return g.ErrNoRecord
	}

	err := app.checkPlatform(pubkey)
	if err == nil {
		return nil
	}
	official,err := app.getOfficial(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, official.Pubkey) {
		return g.ErrAuthenticate
	}
	if official.Active == false {
		return g.ErrNoRight
	}

	var right int32 = 0
	for _, v := range official.Roleids {
		role ,err := app.getRole(v)
		if err != nil {
			return err
		}
		right |= role.Rights
	}

	if right & int32(g.RightType_RsReviewWithdraw) == 0 {
		return g.ErrNoRight
	}

	return nil
}

func (app *Fzmsupply) checkMortgage(req *g.RequestMortgage,uid string,pubkey []byte) error {
	user,err := app.getUser(uid)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, user.Pubkey) {
		return g.ErrAuthenticate
	}

	m,err := app.getMortgage(uid)
	if err == nil && m != nil{
		return g.ErrMortgageExist
	}
	if req.MortgageAmount <= 0 {
		return g.ErrIllegalAmount
	}

	return nil
}

func (app *Fzmsupply) checkConfirmMortgage (req *g.RequestConfirmMortgage,uid string,pubkey []byte) error {

	if req.Rmb < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkRedeem (req *g.RequestRedeem,uid string,pubkey []byte) error {

	if req.RepayAmount < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkCancelMortgage (req *g.RequestCancelMortgage,uid string,pubkey []byte) error {

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkReceive (req *g.RequestReceive,uid string,pubkey []byte) error {

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkSell (req *g.RequestSell,uid string,pubkey []byte) error {

	if req.Amount < 0 || req.FinanceAmount < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkCancel (req *g.RequestCancel,uid string,pubkey []byte) error {

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkInvestAsset (req *g.RequestInvestAsset,uid string,pubkey []byte) error {

	if req.PayAmount < 0 || req.AssetAmount < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkPay (req *g.RequestPay,uid string,pubkey []byte) error {

	if req.Amount < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkRecover (req *g.RequestRecover,uid string,pubkey []byte) error {
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkCancelPay (req *g.RequestCancelPay,uid string,pubkey []byte) error {

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkBatchCash (req *g.RequestBatchCash,uid string,pubkey []byte) error {

	if len(req.Owners) == 0 {
		return g.ErrNoRecord
	}
	if req.Amount < 0 {
		return g.ErrIllegalAmount
	}
	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) checkOfflineCash (req *g.RequestOfflineCash,uid string,pubkey []byte) error {

	if req.PayAmount < 0 {
		return g.ErrIllegalAmount
	}
	if req.AssetId == "" || req.PaymentId == "" {
		return g.ErrWrongValue
	}

	return app.checkIdentity(pubkey, uid)
}

func (app *Fzmsupply) getPlatform() ( platform *g.Platform,err error ) {
	var p g.Platform
	key := fzmsupplyKeyPlatform()
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&p)
	if err != nil {
		return nil,err
	}
	return &p,nil
}


func (app *Fzmsupply) getOfficial(id string) ( official *g.Official,err error ) {
	var p g.Official
	key := fzmsupplyKeyOfficial(id)
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&p)
	if err != nil {
		return nil,err
	}
	return &p,nil
}

func (app *Fzmsupply) getUser(id string) ( u *g.User,err error ) {
	var user g.User
	key := fzmsupplyKeyUser(id)
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&user)
	if err != nil {
		return nil,err
	}
	return &user,nil
}

func (app *Fzmsupply) getRole(id string) ( u *g.Role,err error ) {
	var role g.Role
	key := fzmsupplyKeyRole(id)
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&role)
	if err != nil {
		return nil,err
	}
	return &role,nil
}


func (app *Fzmsupply) getWithdraw(id string) ( u *g.Withdraw,err error ) {
	var with g.Withdraw
	key := fzmsupplyKeyWithdraw(id)
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&with)
	if err != nil {
		return nil,err
	}
	return &with,nil
}


func (app *Fzmsupply) getMortgage(id string) ( u *g.Mortgage,err error ) {
	var m g.Mortgage
	key := fzmsupplyKeyWithdraw(id)
	value,err := app.GetStateDB().Get(key)
	if err != nil {
		return nil,err
	}
	err = types.Decode(value,&m)
	if err != nil {
		return nil,err
	}
	return &m,nil
}