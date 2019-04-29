package types

import (
	"github.com/33cn/chain33/types"

	"reflect"
)

var PreSaleX = "presale"

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(PreSaleX))
	types.RegistorExecutor(PreSaleX, NewType())
	types.RegisterDappFork(PreSaleX, "Enable", 0)
}

type PreSaleType struct {
	types.ExecTypeBase
}

func NewType() *PreSaleType {
	c := &PreSaleType{}
	c.SetChild(c)
	return c
}

func (presale *PreSaleType) GetPayload() types.Message {
	return &PreSale{}
}

func (presale *PreSaleType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"MarketReg":			PreSaleActionMarketReg,
		"AdminReg": 			PreSaleActionAdminReg,
		"UserReg":       		PreSaleActionUserReg,
		"MarketRec":    		PreSaleActionMarketRec,
		"UserRec":      		PreSaleActionUserRec,
		"AddPdInfo":     		PreSaleActionAddPdInfo,
		"OnOffPd":      		PreSaleActionOnOffPd,
		"UserPay": 				PreSaleActionUserPay,
		"Address":    		PreSaleActionAddAddress,
		"OrderTrans":     		PreSaleActionOrderTrans,
		"MarketShip":   		PreSaleActionMarketShip,
		"Logistics":			PreSaleActionLogistics,
		"UserRecipt":			PreSaleActionUserRecipt,
		"OverDue" : 			PreSaleActionOverDue,
		"WithDraw":				PreSaleActionWithDraw,
		"ChangePd" :			PreSaleActionChangePd,
		"UserTokenWithDraw": 	PreSaleActionUserTokenWithDraw,
		"UserShip" : 			PreSaleActionUserShip,
	}
}

func (at *PreSaleType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogPreSaleMarketReg: {reflect.TypeOf(MarketInfoPreSale{}), "TyLogPreSaleMarketReg"},
		TyLogPreSaleAdminReg: {reflect.TypeOf(AdminInfoPreSale{}), "TyLogPreSaleAdminReg"},
		TyLogPreSaleUserReg: {reflect.TypeOf(UserInfoSale{}), "TyLogPreSaleUserReg"},
		TyLogPreSaleMarketRec: {reflect.TypeOf(MarketInfoPreSale{}), "TyLogPreSaleMarketRec"},
		TyLogPreSaleUserRec: {reflect.TypeOf(UserInfoSale{}), "TyLogPreSaleUserRec"},
		TyLogPreSaleAddPdInfo: {reflect.TypeOf(PdInfoPresale{}), "TyLogPreSaleAddPdInfo"},
		TyLogPreSaleOnOffPd: {reflect.TypeOf(PdInfoPresale{}), "TyLogPreSaleOnOffPd"},
		TyLogPreSaleUserPay: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleUserPay"},
		TyLogPreSaleAddAddress: {reflect.TypeOf(UserInfoSale{}), "TyLogPreSaleAddAddress"},
		TyLogPreSaleOrderTrans: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleOrderTrans"},
		TyLogPreSaleMarketShip: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleMarketShip"},
		TyLogPreSaleLogistics: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleLogistics"},
		TyLogPreSaleUserRecipt: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleUserRecipt"},
		TyLogPreSaleOverDue: {reflect.TypeOf(OverDuePreSale{}), "TyLogPreSaleOverDue"},
		TyLogPreSaleWithDraw: {reflect.TypeOf(WithDrawPS{}), "TyLogPreSaleWithDraw"},
		TyLogPreSaleChangePd: {reflect.TypeOf(PdInfoPresale{}), "TyLogPreSaleChangePd"},
		TyLogPreSaleUserTokenWithDraw: {reflect.TypeOf(UserInfoSale{}), "TyLogPreSaleUserTokenWithDraw"},
		TyLogPreSaleUserShip: {reflect.TypeOf(OrderInfoPreSale{}), "TyLogPreSaleUserShip"},
	}
}
