package executor

import (
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	//"github.com/33cn/chain33/common/address"
	"encoding/hex"
	pty "github.com/33cn/plugin/plugin/dapp/presale/types"
	"github.com/33cn/chain33/system/dapp"
	"bytes"
	"github.com/pkg/errors"
	"sort"
)

type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	index        int
}

func NewAction(t *PreSale, tx *types.Transaction,index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{t.GetCoinsAccount(), t.GetStateDB(), hash, fromaddr,
		t.GetBlockTime(), t.GetHeight(),dapp.ExecAddress(string(tx.Execer)),index}
}

//judge+id 与 pubkey绑定 商家加market
func (a *Action) execMarketReg(marketReg *pty.PreSaleMarketReg, pubkey []byte) (*types.Receipt, error) {
	marketName := marketReg.MarketName
	marketId := marketReg.MarketId
	contract := marketReg.Contact
	phoneNum := marketReg.PhoneNum
	// btyAddress := marketReg.BtyAddress
	// ethAddress := marketReg.EthAddress
	currency := marketReg.Currency

	keyJudge := proKeyMarketId(marketId)
	keyMarketPub := proKeyMarketIdJudge(marketId)

	_, err := a.db.Get(keyMarketPub)
	if err == nil {
		return nil, pty.ErrorUserExist
	}

	if marketId == 0 || marketName == "" {
		clog.Error("presale marketReg", "hash", a.txhash, "something is empty", pubkey, marketName)
		return nil, pty.ErrorWrongMessage
	}

	var marketInfo pty.PreSaleMarketInfo //下面为赋值操作
	marketInfo.MarketName = marketName
	marketInfo.MarketId = marketId
	marketInfo.Contact = contract
	marketInfo.PhoneNum = phoneNum
	marketInfo.IsTokenSup = marketReg.IsTokenSup
	// currencyEth := &pty.PreSaleCurrency{}
	// currencyEth.Address = ethAddress
	// currencyEth.UseAble = 0
	// currencyEth.Frozen = 0
	// currencyBty := &pty.PreSaleCurrency{}
	// currencyBty.Address = btyAddress
	// currencyBty.UseAble = 0
	// currencyBty.Frozen = 0
	for _, info := range currency {
		marketCurrency := &pty.PreSaleCurrency{}
		marketCurrency.Address = info.Address
		marketCurrency.Frozen = 0
		marketCurrency.UseAble = 0
		marketCurrency.CurrencyName = info.Name

		marketInfo.Currency = append(marketInfo.Currency, marketCurrency)
	}

	operate := &pty.PreSaleOpreation{}
	operate.Status = "商家注册"

	marketInfo.Operate = operate

	tokenInfo := &pty.PreSaleToken{}
	tokenInfo.MarketId = marketId
	tokenInfo.UseAble = 0
	tokenInfo.TokenFrozen = 0
	marketInfo.Token = tokenInfo

	// marketInfo.Bty = currencyBty
	// marketInfo.Eth = currencyEth

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&marketInfo)
	err = a.execSetdb(keyJudge, value)
	if err != nil {
		clog.Error("presale marketReg", "hash", a.txhash, "set marketdb failed", err)
		return nil, err
	}
	err = a.db.Set(keyMarketPub, pubkey)
	if err != nil {
		return nil, err
	}
	// err= a.execSetdb(pubkey, keyJudge)
	// if err != nil {
	// 	clog.Error("presale marketReg", "hash", a.txhash, "set pubkeydb failed", err2)
	// 	return nil, err
	// }
	//上链
	kv = append(kv, &types.KeyValue{Key:key, Value:value}, &types.KeyValue{Key:keyMarketPub, Value:pubkey}, &types.KeyValue{Key:keyJudge, Value:value})
	//log := &pty.MarketInfoPreSale{Order: &order}
	//logs = append(logs, &types.ReceiptLog{Ty: TyLogPreSaleMarketReg, Log: types.Encode(log)})
	log := &pty.MarketInfoPreSale{MarketInfo: &marketInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleMarketReg, Log: types.Encode(log)})
	//receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execSetdb(key []byte, value []byte) error {
	err := a.db.Set(key, value)
	if err != nil {
		return err
	}
	return nil
}

func (a *Action) execAdminReg(adminReg *pty.PreSaleAdminReg, pubkey []byte) (*types.Receipt, error) {
	adminName := adminReg.AdminName
	adminId := adminReg.AdminId
	if adminId == 0 || adminName == "" {
		clog.Error("presale adminReg", "something is empty", adminId, adminName)
		return nil, pty.ErrorWrongMessage
	}

	keyAdminReg := proKeyUserId(adminId)
	keyAdminPub := proKeyUidJudge(adminId)
	_, err := a.db.Get(keyAdminPub)
	if err == nil {
		return nil, pty.ErrorUserExist
	}

	var order pty.PreSaleAdminInfo
	order.AdminName = adminName
	order.AdminId = adminId
	order.Status = "工作人员注册"

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	err1 := a.execSetdb(keyAdminPub, pubkey)
	if err1 != nil {
		clog.Error("presale adminReg", "hash", a.txhash, "set pubkey db failed", err1)
		return nil, err1
	}
	err2 := a.execSetdb(keyAdminReg, value)
	if err2 != nil {
		clog.Error("presale adminReg", "hash", a.txhash, "set admin db failed", err2)
		return nil, err2
	}
	//上链
	kv = append(kv, &types.KeyValue{Key:key, Value:value}, &types.KeyValue{Key:keyAdminPub, Value:pubkey}, &types.KeyValue{Key:keyAdminReg, Value:value})
	log := &pty.AdminInfoPreSale{AdminInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleAdminReg, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserReg(userReg *pty.PreSaleUserReg, pubkey []byte) (*types.Receipt, error) {
	userName := userReg.Name
	userEmail := userReg.Email
	userPhoneNum := userReg.PhoneNum
	userId := userReg.Uid
	// btyAddress := userReg.BtyAddress
	// ethAddress := userReg.EthAddress
	currency := userReg.Currency
	uuid := userReg.Uuid

	keyUserReg := proKeyUserId(userId)
	keyUserPub := proKeyUidJudge(userId)
	// userPubkey, err := a.db.Get([]byte(keyUserReg))
	// if err == nil {
	// 	clog.Error("userReg ", "get pubkey success", userPubkey)
	// 	return nil, pty.ErrorUserExist
	// }

	if userName == "" || userId == 0 {
		clog.Error("presale userreg", "something is empty", userName, userId)
		return nil, pty.ErrorWrongMessage
	}

	_, err := a.db.Get(keyUserPub)
	if err == nil {
		clog.Error("presale userreg", "db is exist", userId)
		return nil, pty.ErrorUserExist
	}
	var userInfo pty.PreSaleUserInfo
	userInfo.Name = userName
	userInfo.Email = userEmail
	userInfo.PhoneNum = userPhoneNum
	userInfo.Uid = userId
	userInfo.Uuid = uuid
	// currencyEth := &pty.PreSaleCurrency{}
	// currencyEth.Address = ethAddress
	// currencyEth.UseAble = 0
	// currencyEth.Frozen = 0

	operate := &pty.PreSaleOpreation{}
	operate.Status = "用户注册"

	userInfo.Operate = operate

	// currencyBty := &pty.PreSaleCurrency{}
	// currencyBty.Address = btyAddress
	// currencyBty.UseAble = 0
	// currencyBty.Frozen = 0

	// userInfo.Bty = currencyBty
	// userInfo.Eth = currencyEth
	for _, info := range currency {
		userCurrency := &pty.PreSaleCurrency{}
		userCurrency.Address = info.Address
		userCurrency.Frozen = 0
		userCurrency.UseAble = 0
		userCurrency.CurrencyName = info.Name

		userInfo.Currency = append(userInfo.Currency, userCurrency)
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&userInfo)
	err1 := a.execSetdb(keyUserPub, pubkey)
	if err1 != nil {
		clog.Error("presale userreg", "hash", a.txhash, "set pubkey db failed", err1)
		return nil, err1
	}
	//	_, err11 := a.db.Get(keyUserPub)
	//clog.Error("presale userreg", "hash", a.txhash, "pubkey", pubkey, "uid is ", value1, "err", err11)

	err2 := a.execSetdb(keyUserReg, value)
	if err2 != nil {
		clog.Error("presale adminReg", "hash", a.txhash, "set user db failed", err2)
		return nil, err2
	}
	//上链
	kv = append(kv, &types.KeyValue{Key:keyUserReg, Value:value}, &types.KeyValue{Key:keyUserPub, Value:pubkey},
		&types.KeyValue{Key:[]byte(key), Value:value})

	log := &pty.UserInfoSale{UserInfo: &userInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserReg, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execMarketRec(marketRec *pty.PreSaleMarketRec, pubkey []byte) (*types.Receipt, error) {
	marketId := marketRec.MarketId
	amount := marketRec.Amount
	currency := marketRec.Currency
	isToken := marketRec.IsToken
	time := marketRec.Time
	toAddress := marketRec.ToAddress
	keyMarketId := proKeyMarketId(marketId)

	err := a.JudgeMarkeyId(marketId, pubkey)
	if err != nil {
		return nil, err
	}

	marketInfo, err := a.getMarketInfo(marketId)
	if err != nil {
		return nil, err
	}

	if amount <= 0 {
		return nil, pty.ErrorRecAmount
	}
	//判断商家id是否匹配
	// if marketInfo.MarketId != marketId {
	// 	clog.Error("presale marketrec", "order.MarketId",
	// 		marketInfo.MarketId, "marketid ", marketId)
	// 	return nil, pty.ErrorWrongMarket
	// }
	//clog.Error("presale marketrec", "order.MarketId",
	//order.MarketId, "marketid ", marketId)
	//先判断是不是token 对资产加加减减
	// if !isToken {
	// 	switch currency {
	// 	case "bty":
	// 		if strings.Compare(toAddress, marketInfo.Bty.Address) != 0 {
	// 			return nil, pty.ErrorAddress
	// 		}
	// 		marketInfo.Bty.UseAble += amount
	// 	case "eth":
	// 		if strings.Compare(toAddress, marketInfo.Eth.Address) != 0 {
	// 			return nil, pty.ErrorAddress
	// 		}
	// 		marketInfo.Eth.UseAble += amount
	// 	default:
	// 		return nil, pty.ErrorWrongCoinType
	// 	}
	// } else {
	// 	switch currency {
	// 	case "bty":
	// 		if strings.Compare(toAddress, marketInfo.Bty.Address) != 0 {
	// 			return nil, pty.ErrorAddress
	// 		}
	// 		marketInfo.Token.UseAble += amount
	// 	case "eth":
	// 		if strings.Compare(toAddress, marketInfo.Eth.Address) != 0 {
	// 			return nil, pty.ErrorAddress
	// 		}
	// 		marketInfo.Token.UseAble += amount
	// 	default:
	// 		return nil, pty.ErrorWrongCoinType
	// 	}
	// }
	c := 0
	if !isToken {
		for _, info := range marketInfo.Currency {
			if currency == info.CurrencyName {
				// if strings.Compare(toAddress, info.Address) == 0 {
				// 	info.UseAble += amount
				// 	c++
				// }
				info.UseAble += amount
				c++

			}
		}
	} else {
		for _, info := range marketInfo.Currency {
			// if strings.Compare(toAddress, info.Address) == 0 && currency == info.CurrencyName {
			if currency == info.CurrencyName {
				marketInfo.Token.UseAble += amount
				c++
			}
		}
	}
	if c == 0 {
		return nil, pty.ErrorRecInfo
	}
	operate := &pty.PreSaleOpreation{}
	operate.Time = time
	operate.Status = "商家充值"
	operate.ToAddress = toAddress

	marketInfo.Operate = operate

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	value := types.Encode(marketInfo)
	key := calcOrderKey(string(a.txhash))
	err = a.execSetdb(keyMarketId, value)
	if err != nil {
		clog.Error("presale marketrec", "hash", a.txhash, "set pubkey db failed", err)
		return nil, err
	}
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value})
	kv = append(kv, &types.KeyValue{Key:keyMarketId, Value:value})
	log := &pty.MarketInfoPreSale{MarketInfo: marketInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleMarketRec, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserRec(userRec *pty.PreSaleUserRec, pubkey []byte) (*types.Receipt, error) {
	uid := userRec.Uid
	amount := userRec.Amount
	time := userRec.Time
	currency := userRec.Currency
	toAddress := userRec.ToAddress

	keyJudge := proKeyUserId(uid)
	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, pty.ErrorRecAmount
	}

	var order pty.PreSaleUserInfo
	userInfo, err3 := a.db.Get(keyJudge)
	if err3 != nil {
		clog.Error("presale userRec", "hash", a.txhash, " get db failed", uid)
		return nil, pty.ErrorUser
	}
	err1 := types.Decode(userInfo, &order)
	if err1 != nil {
		clog.Error("presale UserRec", "hash", a.txhash, "decode order failed",
			uid, "err", err1)
		return nil, err1
	}
	// switch currency {
	// case "bty":
	// 	if strings.Compare(toAddress, order.Bty.Address) != 0 {
	// 		return nil, pty.ErrorAddress
	// 	}
	// 	order.Bty.UseAble += amount
	// case "eth":
	// 	if strings.Compare(toAddress, order.Eth.Address) != 0 {
	// 		return nil, pty.ErrorAddress
	// 	}
	// 	order.Eth.UseAble += amount
	// default:
	// 	return nil, pty.ErrorToken
	// }
	c := 0
	for _, info := range order.Currency {
		if info.CurrencyName == currency {
			// if strings.Compare(info.Address, toAddress) == 0 {
			// 	info.UseAble += amount
			// 	c++
			// }
			info.UseAble += amount
			c++
		}
	}
	if c == 0 {
		return nil, pty.ErrorRecInfo
	}
	operate := &pty.PreSaleOpreation{}
	operate.Time = time
	operate.ToAddress = toAddress
	operate.Status = "用户充值"

	order.Operate = operate

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	err4 := a.execSetdb(keyJudge, value)
	if err4 != nil {
		clog.Error("presale marketrec", "hash", a.txhash, "set pubkey db failed", err3)
		return nil, err4
	}
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyJudge, Value:value})
	log := &pty.UserInfoSale{UserInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserRec, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execAddPdInfo(addPdInfo *pty.PreSaleAddPdInfo, pubkey []byte) (*types.Receipt, error) {
	pdName := addPdInfo.PdName
	pdId := addPdInfo.PdId
	price := addPdInfo.Price
	deposit := addPdInfo.Deposit
	psStart := addPdInfo.PsStart
	psEnd := addPdInfo.PsEnd
	puStart := addPdInfo.PuStart
	puEnd := addPdInfo.PuEnd
	marketId := addPdInfo.MarketId
	marketName := addPdInfo.MarketName
	adminId := addPdInfo.AdminId
	status := addPdInfo.Status

	// keyJudge := proKeyMarketId(marketId)
	// keyMarketPub := proKeyMarketIdJudge(marketId)

	keyPdId := proKeyPdId(pdId)
	// err := a.JudgeUid(adminId, pubkey)
	// if err != nil {
	// 	return nil, err
	// }

	//keyMarketId := proKeyIdReg(marketId) //id:pubkey
	markeInfo, err1 := a.getMarketInfo(marketId)
	if err1 != nil {
		clog.Error("presale addInfo", "hash", a.txhash, "get marketInfo failed",
			marketId)
		return nil, err1
	}

	if markeInfo.MarketId != marketId || markeInfo.MarketName != marketName {
		clog.Error("presale addInfo", "hash", a.txhash, "Info not equal",
			marketId, marketName)
		return nil, pty.ErrorMarketInfo
	}
	_, err := a.db.Get(keyPdId)
	if err == nil {
		return nil, pty.ErrorPdExist
	}
	var orderSet pty.PresalePdInfo
	orderSet.PdName = pdName
	orderSet.PdId = pdId
	orderSet.Price = price
	orderSet.Deposit = deposit
	orderSet.PsStart = psStart
	orderSet.PsEnd = psEnd
	orderSet.PuStart = puStart
	orderSet.PuEnd = puEnd
	orderSet.MarketId = marketId
	orderSet.MarketName = marketName
	orderSet.AdminId = adminId
	orderSet.Status = status

	if status == "" {
		return nil, pty.ErrorStatus
	}
	if status == "已删除" {
		return nil, pty.ErrorPdDelete
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&orderSet)
	err3 := a.execSetdb(keyPdId, value)
	if err3 != nil {
		clog.Error("presale addInfo", "hash", a.txhash, "set pubkey db failed", err3)
		return nil, err3
	}
	//上链
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyPdId, Value:value})
	log := &pty.PdInfoPresale{PdInfo: &orderSet}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleAddPdInfo, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execGetOnOffPd(onOffPd *pty.PreSaleOnOffPd, pubkey []byte) (*types.Receipt, error) {
	pdId := onOffPd.PdId
	status := onOffPd.Status
	adminId := onOffPd.AdminId
	//time := onOffPd.Time
	keyJudgePd := proKeyPdId(pdId)

	err := a.JudgeMarkeyId(adminId, pubkey)
	if err != nil {
		return nil, err
	}
	//marketid和marketname是否一致 todo
	//两种价格上建立绑定

	var order pty.PresalePdInfo
	orderInfo, err1 := a.db.Get(keyJudgePd)
	if err1 != nil {
		clog.Error("presale onOffPd", "hash", a.txhash, "get kv failed", pdId)
		return nil, pty.ErrorWrongPd
	}
	err = types.Decode(orderInfo, &order)
	if err != nil {
		clog.Error("presale onOffPd", "hash", a.txhash, "decode failed", pdId)
		return nil, err
	}
	//if order.PsEnd < time {
	//	return nil, pty.ErrorPdStatus
	//}
	if order.Status == "已删除" {
		return nil, pty.ErrorPdDelete
	}
	if order.Status == "已上线" {
		clog.Error("presale onOffPd", "already onSale", pdId)
		if status == true {
			return nil, pty.ErrorPdInfo
		} else {
			order.Status = "下线"
		}
	} else if order.Status == "未上线" {
		if status == true {
			order.Status = "已上线"
		} else {
			order.Status = "下线"
		}
	} else if order.Status == "下线" {
		if status == true {
			order.Status = "已上线"
		} else {
			return nil, pty.ErrorPdInfoOff
		}
	} else {
		return nil, pty.ErrorWrongMessage
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	//上链
	err3 := a.execSetdb(keyJudgePd, value)
	if err3 != nil {
		clog.Error("presale addInfo", "hash", a.txhash, "set pubkey db failed", err3)
		return nil, err3
	}
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyJudgePd, Value:value})
	log := &pty.PdInfoPresale{PdInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleOnOffPd, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserPay(userPay *pty.PreSaleUserPay, pubkey []byte) (*types.Receipt, error) {
	//order,err := a.getOrderInfo(userPay.OrderId)
	//if err != nil {
	//	return nil,err
	//}
	//if order == nil {
	//	return nil,pty.ErrorOrderNotExist
	//}
	//if userPay.MarketId != order.MarketId {
	//	return nil,pty.ErrorWrongMarket
	//}
	market, err := a.getMarketInfo(userPay.MarketId)
	if err != nil {
		return nil, err
	}
	if market == nil {
		return nil, pty.ErrorMarketNotExist
	}
	payId := userPay.PayId
	switch payId {
	case pty.PreSalePayType_PreSalePayBook:
		return a.execUserBook(userPay, pubkey, market.IsTokenSup)
	case pty.PreSalePayType_PreSalePayLast:
		return a.execUserPayLast(userPay, pubkey, market.IsTokenSup)
	default:
		return nil, pty.ErrorWrongMessage
	}

	//先判断商家的marketid是否和orderid一致
}

func (a *Action) getOrderInfo(orderId int64) (*pty.PreSaleOrderInfo, error) {
	var orderInfo pty.PreSaleOrderInfo
	keyOrder := proKeyOrderId(orderId)
	value, err := a.db.Get(keyOrder)
	if err != nil {
		return &orderInfo, err
	}
	err = types.Decode(value, &orderInfo)
	if err != nil {
		return &orderInfo, err
	}
	return &orderInfo, nil
}

func (a *Action) execUserBook(userBook *pty.PreSaleUserPay, pubkey []byte, isTokenSup bool) (*types.Receipt, error) {
	orderId := userBook.OrderId
	pdId := userBook.PdId
	uid := userBook.Uid
	currency := userBook.CurrencyName
	tokenAmount := userBook.TokenAmount
	currencyAmount := userBook.CurrencyAmount
	time := userBook.Time
	marketId := userBook.MarketId
	pdAmount := userBook.Amount

	//pd有没有数量 如果有那就是整数倍
	if tokenAmount < 0 {
		return nil, pty.ErrorTokenAmount
	}

	var order pty.PresalePdInfo
	//var orderUser pty.PreSaleUserInfo
	//var orderMarket pty.PreSaleMarketInfo
	var orderSet pty.PreSaleOrderInfo
	//keyUId := proKeyIdReg(uid)

	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}

	keyPdId := proKeyPdId(pdId)
	pdInfo, err := a.db.Get(keyPdId)
	if err != nil {
		clog.Error("presale userBook", "hash", a.txhash, "get kvdb failed",
			pdId, "err", err)
		return nil, err
	}
	err1 := types.Decode(pdInfo, &order)
	if err1 != nil {
		clog.Error("presale userBook", "hash", a.txhash, "decode order failed",
			pdId, "err", err1)
		return nil, err1
	}
	if order.Status == "已删除" {
		return nil, pty.ErrorPdDelete
	}
	//先判断订单状态对不对
	if order.Status == "下线" || order.Status == "未上线" {
		clog.Error("presale userBook", "hash", a.txhash, "wrong message", userBook)
		return nil, pty.ErrorWrongPd
	}
	//if order.PsEnd < time {
	//	return nil, ErrorBookTime
	//}
	// if order.PdId != pdId {
	// 	clog.Error("presale userBook", "hash", a.txhash, "wrong message", userBook)
	// 	return nil, pty.ErrorWrongMessage
	// }
	//转让的token数量要等于数量乘以定价
	//if tokenAmount != pdAmount*order.Deposit {
	//	return nil, pty.ErrorTokenTrans
	//}
	//更改用户信息
	keyUserReg := proKeyUserId(uid)
	marketInfo, err3 := a.getMarketInfo(marketId)
	if err3 != nil {
		return nil, err3
	}
	userInfo, err2 := a.getUserInfo(uid)
	if err2 != nil {
		return nil, err2
	}
	keyOrderId := proKeyOrderId(orderId)
	_, err = a.db.Get(keyOrderId)
	if err == nil {
		return nil, pty.ErrorOrderExist
	}
	keyMaketId := proKeyMarketId(marketId)
	//token支付
	if isTokenSup {
		//添加用户token账户
		b := 0
		//用户冻结token增加
		for _, res := range userInfo.Token {
			if res.MarketId == marketId {
				res.TokenFrozen += tokenAmount
				b++
			}
		}
		if b == 0 {
			tokenInfo := &pty.PreSaleToken{}
			tokenInfo.MarketId = marketId
			tokenInfo.UseAble = 0
			tokenInfo.TokenFrozen = tokenAmount
			userInfo.Token = append(userInfo.Token, tokenInfo)
		}
		//判断商家账户够不够token

		if marketInfo.Token.UseAble < tokenAmount {
			return nil, pty.ErrorMarketOwn
		}
		//1.商家可用Token减少
		marketInfo.Token.UseAble -= tokenAmount
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	orderSet.OrderId = orderId
	orderSet.PdId = pdId
	orderSet.Uid = uid
	orderSet.MarketId = marketId
	orderSet.TokenAmount = tokenAmount
	orderSet.UserStatus = "已预订"
	orderSet.OperateTime = time
	orderSet.Status = "用户预订"
	orderSet.PdAmount = pdAmount
	orderSet.PsStart = order.PsStart
	orderSet.PsEnd = order.PsEnd
	orderSet.PuStart = order.PuStart
	orderSet.PuEnd = order.PuEnd
	book := &pty.PreSaleAmount{}
	book.CurrencyName = currency
	book.Amount = currencyAmount
	orderSet.Book = book
	// if currency == "bty" {
	// 	userInfo.Bty.UseAble -= currencyAmount
	// 	marketInfo.Bty.Frozen += currencyAmount
	// 	orderSet.BtyAmount += currencyAmount
	// } else if currency == "eth" {
	// 	userInfo.Eth.UseAble -= currencyAmount
	// 	marketInfo.Eth.Frozen += currencyAmount
	// 	orderSet.EthAmount += currencyAmount
	// } else {
	// 	return nil, pty.ErrorToken
	// }
	//对虚拟币进行操作
	c := 0
	for _, info := range userInfo.Currency {
		if info.CurrencyName == currency {
			info.UseAble -= currencyAmount
			c++
		}
	}
	d := 0
	for _, infomar := range marketInfo.Currency {
		if infomar.CurrencyName == currency {
			infomar.Frozen += currencyAmount
			d++
		}
	}

	if c == 0 || d == 0 {
		return nil, pty.ErrorRecInfo
	}
	valueUser := types.Encode(userInfo)
	err = a.db.Set(keyUserReg, valueUser)
	if err != nil {
		clog.Error("presale userBook", "hash", a.txhash, "set userdb failed",
			keyUserReg, "err", err)
		return nil, err
	}
	valueOrderId := types.Encode(&orderSet)
	err = a.db.Set(keyOrderId, valueOrderId)
	if err != nil {
		clog.Error("presale userBook", "hash", a.txhash, "set userdb failed",
			marketId, "err", err)
		return nil, err
	}
	valueMarket := types.Encode(marketInfo)
	err = a.db.Set(keyMaketId, valueMarket)
	if err != nil {
		clog.Error("presale userBook", "hash", a.txhash, "set userdb failed",
			marketInfo, "err", err)
		return nil, err
	}
	key := calcOrderKey(string(a.txhash))
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:valueOrderId}, &types.KeyValue{Key:keyUserReg, Value:valueUser},
		&types.KeyValue{Key:keyMaketId, Value:valueMarket}, &types.KeyValue{Key:keyOrderId, Value:valueOrderId})
	log := &pty.OrderInfoPreSale{OrderInfo: &orderSet}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserPay, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserPayLast(userPayLast *pty.PreSaleUserPay, pubkey []byte, isTokenSup bool) (*types.Receipt, error) {
	orderId := userPayLast.OrderId
	uid := userPayLast.Uid
	// currency := ""
	// switch userPayLast.CurrencyId {
	// case 1:
	// 	currency = "bty"
	// case 2:
	// 	currency = "eth"
	// default:
	// 	return nil, ErrorCurrencyId
	// }
	currency := userPayLast.CurrencyName
	tokenAmount := userPayLast.TokenAmount
	currencyAmount := userPayLast.CurrencyAmount
	time := userPayLast.Time
	marketId := userPayLast.MarketId
	addressId := userPayLast.AddressId
	pdId := userPayLast.PdId

	// var orderSet pty.PreSaleOrderInfo
	var pdInfo pty.PresalePdInfo

	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}
	keyPdInfo := proKeyPdId(pdId)
	keyOrderId := proKeyOrderId(orderId)
	orderSet, err1 := a.getOrderInfo(orderId)
	if err1 != nil {
		clog.Error("presale userPayLast", "hash", a.txhash, "get orderInfo failed",
			orderId)
		return nil, err1
	}
	pdOrder, err := a.db.Get(keyPdInfo)
	if err != nil {
		return nil, err
	}
	_ = types.Decode(pdOrder, &pdInfo)
	if pdInfo.Status == "已删除" {
		return nil, pty.ErrorPdDelete
	}
	//order信息中orderid 和 Uid 以及marketid是否对应
	if orderSet.OrderId != orderId || orderSet.Uid != uid ||
		orderSet.MarketId != marketId {
		return nil, pty.ErrorOrderMessage
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	orderSet.TokenAmount += tokenAmount
	orderSet.UserStatus = "待发货"
	orderSet.OperateTime = time
	orderSet.Uid = uid
	orderSet.Status = "用户付尾款"
	payLast := &pty.PreSaleAmount{}
	payLast.Amount = currencyAmount
	payLast.CurrencyName = currency
	orderSet.LastPay = payLast
	//if orderSet.PuEnd < time {
	//	return nil, ErrorLastPayTime
	//}
	keyMarketId := proKeyMarketId(marketId)
	orderMarket, err4 := a.getMarketInfo(marketId)
	if err4 != nil {
		clog.Error("presale userBook", "hash", a.txhash, "get market info failed",
			marketId, "err", err4)
		return nil, err4
	}
	if orderMarket.IsTokenSup {
		if orderMarket.Token.UseAble < tokenAmount {
			return nil, pty.ErrorMarketOwn
		}
	}

	//改变商家资产情况
	c := 0
	for _, info := range orderMarket.Currency {
		if info.CurrencyName == currency {
			info.Frozen += currencyAmount
			if isTokenSup {
				orderMarket.Token.UseAble -= tokenAmount
			}

			c++
		}
	}
	if c == 0 {
		return nil, pty.ErrorRecInfo
	}
	keyUserReg := proKeyUserId(uid)
	orderUser, err3 := a.getUserInfo(uid)
	if err3 != nil {
		clog.Error("presale userPayLast", "hash", a.txhash, "set uid info failed",
			uid, "err", err3)
		return nil, err3
	}

	//用户的账户余额改变
	//orderUser.Token = orderMarket.TokenName
	if isTokenSup {
		for _, res := range orderUser.Token {
			if res.MarketId == marketId {
				res.TokenFrozen += tokenAmount
			}
		}
	}

	for _, info := range orderUser.Currency {
		if info.CurrencyName == currency {
			info.UseAble -= currencyAmount
			// orderUser.Token.Frozen += tokenAmount
		}
	}
	valueUser := types.Encode(orderUser)
	err = a.db.Set(keyUserReg, valueUser)
	if err != nil {
		clog.Error("presale userPayLast", "hash", a.txhash, "set user kvdb failed",
			keyUserReg, "err", err)
		return nil, err
	}
	d := 0
	for _, address := range orderUser.Address {
		if address.AddressId == addressId {
			orderSet.Address = address.Address
			orderSet.UserName = address.UserName
			orderSet.PhoneNum = address.PhoneNum
			d++
		}
	}
	//if orderSet.TokenAmount != pdInfo.Price*orderSet.PdAmount {
	//	return nil, pty.ErrorTokenAmount
	//}
	if d == 0 {
		//return nil, pty.ErrorAddressId
	}
	valueOrderId := types.Encode(orderSet)
	err8 := a.db.Set(keyOrderId, valueOrderId)
	if err8 != nil {
		clog.Error("presale userBook", "hash", a.txhash, "set orderdb failed",
			orderId, "err", err8)
		return nil, err8
	}
	valueMarket := types.Encode(orderMarket)
	err = a.db.Set(keyMarketId, valueMarket)
	if err != nil {
		clog.Error("presale userBook", "hash", a.txhash, "set userdb failed",
			marketId, "err", err)
		return nil, err
	}
	key := calcOrderKey(string(a.txhash))

	kv = append(kv, &types.KeyValue{Key:key, Value:valueOrderId}, &types.KeyValue{Key:keyUserReg, Value:valueUser},
		&types.KeyValue{Key:keyMarketId, Value:valueMarket}, &types.KeyValue{Key:keyOrderId,Value: valueOrderId})
	log := &pty.OrderInfoPreSale{OrderInfo: orderSet}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserPay, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execAddAddress(address *pty.PreSaleAddAddress, pubkey []byte) (*types.Receipt, error) {
	userName := address.UserName
	phoneNum := address.PhoneNum
	address1 := address.Address
	uid := address.Uid
	addressId := address.AddressId
	//keyJudge := proKeyIdReg(uid)
	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}
	keyUser := proKeyUserId(uid)
	userInfo, err := a.getUserInfo(uid)
	if err != nil {
		return nil, err
	}
	b := 0
	for _, value := range userInfo.Address {
		if value.AddressId == addressId {
			value.UserName = userName
			value.PhoneNum = phoneNum
			value.Address = address1
			b++
		}
	}
	if b == 0 {
		newAddress := &pty.PreSaleAddress{}
		newAddress.UserName = userName
		newAddress.PhoneNum = phoneNum
		newAddress.Address = address1
		newAddress.AddressId = addressId

		userInfo.Address = append(userInfo.Address, newAddress)
	}
	operation := &pty.PreSaleOpreation{}
	operation.Status = "新增用户地址"

	userInfo.Operate = operation
	// var order pty.PreSaleOrderInfo
	// keyOrderId := proKeyOrderId(orderId)
	// order, err1 := a.getOrderInfo(orderId)
	// if err1 != nil {
	// 	clog.Error("presale addAddress", "hash", a.txhash, "get orderInfo failed",
	// 		orderId, "err", err1)
	// 	return nil, err1
	// }
	// if order.UserStatus != "待发货" || order.MarketStatus != "" || order.Uid != uid {
	// 	return nil, pty.ErrorOrder
	// }
	// if order.Uid != uid {
	// 	return nil, pty.ErrorWrongUid
	// }
	// order.UserStatus = "待收货"
	// order.Address = Addaddress
	// order.PhoneNum = phoneNum
	// order.UserName = userName
	// order.Status = "用户添加收货地址"

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(userInfo)
	kv = append(kv, &types.KeyValue{Key:key, Value:value}, &types.KeyValue{Key:keyUser, Value:value})
	log := &pty.UserInfoSale{UserInfo: userInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleAddAddress, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execOrderTrans(orderTrans *pty.PreSaleOrderTrans, pubkey []byte) (*types.Receipt, error) {
	orderId := orderTrans.OrderId
	toUUid := orderTrans.ToUUid //todo 证明touid存在
	toUid := orderTrans.ToUid
	uid := orderTrans.Uid
	time := orderTrans.Time
	newOrderId := orderTrans.NewOrderId
	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}
	// var order pty.PreSaleOrderInfo
	var orderToUser pty.PreSaleUserInfo
	keyOrderId := proKeyOrderId(newOrderId)
	keyOldOrderId := proKeyOrderId(orderId)

	// keyOldOrderId:=proKeyOrderId(newOrderId)
	order, err1 := a.getOrderInfo(orderId)
	if err1 != nil {
		clog.Error("presale orderTrans", "get order info failed", orderId,
			"err", err1)
		return nil, err1
	}
	if order.Uid != uid {
		return nil, pty.ErrorWrongUid
	}
	marketId := order.MarketId
	marketInfo, err2 := a.getMarketInfo(marketId)
	if err2 != nil {
		clog.Error("presale orderTrans", "get market info failed", marketId,
			"err", err2)
		return nil, err2
	}
	//预订状态下才能转让订单
	//if order.UserStatus != "已预订" {
	//	return nil, pty.ErrorOrder
	//}
	order.OrderId = newOrderId
	order.Status = "用户订单转让"
	order.Uid = toUid
	order.OperateTime = time
	//if order.PsEnd < time {
	//	return nil, ErrorBookTime
	//}
	keyUserId := proKeyUserId(uid)
	orderUser, err5 := a.getUserInfo(uid)
	if err5 != nil {
		clog.Error("presale orderTrans", "get user info failed", uid,
			"err", err1)
		return nil, err5
	}
	if marketInfo.IsTokenSup {
		for _, res := range orderUser.Token {
			if res.MarketId == marketId {
				res.TokenFrozen -= order.TokenAmount
				if res.TokenFrozen < 0 {
					return nil, pty.ErrorUserFrozen
				}
			}
		}
	}


	keyToUserId := proKeyUserId(toUid)
	toUserInfo, err2 := a.db.Get(keyToUserId)
	if err2 != nil {
		return nil, pty.ErrorTouser
	}
	err = types.Decode(toUserInfo, &orderToUser)
	if err != nil {
		return nil, err
	}
	b := 0
	if orderToUser.Uuid != toUUid {
		return nil, pty.ErrorUUid
	}
	if marketInfo.IsTokenSup {
		for _, res := range orderToUser.Token {
			if res.MarketId == marketId {
				res.TokenFrozen += order.TokenAmount
				b++
			}
		}
		if b == 0 {
			userToken := &pty.PreSaleToken{}
			userToken.MarketId = marketId
			userToken.TokenFrozen += order.TokenAmount
			orderToUser.Token = append(orderToUser.Token, userToken)
		}
	}


	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(order)
	valueToUser := types.Encode(&orderToUser)
	valueUser := types.Encode(orderUser)
	err4 := a.db.Set(keyOrderId, value)
	if err4 != nil {
		clog.Error("presale orderTrans", "hash", a.txhash, "set db failed",
			keyOrderId, "err", err4)
		return nil, err4
	}
	order.Status = "已转让无法操作"
	valueOld := types.Encode(order)

	err = a.db.Set(keyOldOrderId, valueOld)
	if err != nil {
		return nil, err
	}
	err = a.db.Set(keyToUserId, valueToUser)
	if err != nil {
		return nil, err
	}
	err = a.db.Set(keyUserId, valueUser)
	if err != nil {
		return nil, err
	}
	//上链
	kv = append(kv, &types.KeyValue{Key:key, Value:value}, &types.KeyValue{Key:keyOrderId, Value:value},
		&types.KeyValue{Key:keyToUserId, Value:valueToUser}, &types.KeyValue{Key:keyUserId, Value:valueUser},
		&types.KeyValue{Key:keyOldOrderId, Value:valueOld})
	log := &pty.OrderInfoPreSale{OrderInfo: order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleOrderTrans, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execMarketShip(marketShip *pty.PreSaleMarketShip, pubkey []byte) (*types.Receipt, error) {
	orderId := marketShip.OrderId
	adminId := marketShip.AdminId
	time := marketShip.Time
	express := marketShip.Express
	orderNum := marketShip.OrdNum
	err := a.JudgeMarkeyId(adminId, pubkey)
	if err != nil {
		return nil, err
	}
	var order pty.PreSaleOrderInfo
	keyOrderId := proKeyOrderId(orderId)
	orderInfo, err := a.db.Get(keyOrderId)
	if err != nil {
		clog.Error("presale marketShip", "hash", a.txhash, "get kvdb failed",
			keyOrderId, "err", err)
	}
	err1 := types.Decode(orderInfo, &order)
	if err1 != nil {
		clog.Error("presale marketShip", "hash", a.txhash, "decode order failed",
			orderInfo, "err", err)
		return nil, err
	}
	if order.UserStatus != "用户提货" || order.Address == "" {
		return nil, pty.ErrorOrder
	}

	order.MarketStatus = "已发货"
	order.Status = "商家已发货"
	order.OperateTime = time
	order.OrdNum = orderNum
	order.Express = express
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	//上链
	a.db.Set(keyOrderId, value)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyOrderId, Value:value})
	log := &pty.OrderInfoPreSale{OrderInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleMarketShip, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execLogistics(logistics *pty.PreSaleLogistics, pubkey []byte) (*types.Receipt, error) {
	orderId := logistics.OrderId
	express := logistics.Express
	ordNum := logistics.OrdNum
	time := logistics.Time
	marketId := logistics.Marketid
	err := a.JudgeMarkeyId(marketId, pubkey)
	if err != nil {
		return nil, err
	}
	var order pty.PreSaleOrderInfo
	keyOrderId := proKeyOrderId(orderId)
	orderInfo, err1 := a.db.Get(keyOrderId)
	if err1 != nil {
		clog.Error("presale logistics", "hash", a.txhash, "get kvdb failed",
			keyOrderId, "err", err1)
		return nil, err1
	}
	err = types.Decode(orderInfo, &order)
	if err != nil {
		clog.Error("presale logistics", "hash", a.txhash, "decode order failed",
			orderInfo, "err", err)
		return nil, err
	}
	if order.MarketId != marketId {
		return nil, pty.ErrorMarketOprea
	}
	if order.MarketStatus != "已发货" {
		return nil, pty.ErrorOrder
	}
	order.Express = express
	order.OrdNum = ordNum
	order.ExpressTime = time
	order.OperateTime = time
	order.Status = "商家添加物流信息"

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	//上链
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyOrderId,Value: value})
	log := &pty.OrderInfoPreSale{OrderInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleLogistics, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserRecipt(userRecipt *pty.PreSaleUserRecipt, pubkey []byte) (*types.Receipt, error) {
	orderId := userRecipt.OrderId
	time := userRecipt.Time
	uid := userRecipt.Uid

	m := 0
	adminPubkey, _ := hex.DecodeString(adminPubkey)
	if bytes.Equal(adminPubkey, pubkey) {
		m++
	}
	if m == 0 {
		err := a.JudgeUid(uid, pubkey)
		if err != nil {
			return nil, err
		}
	}
	var order pty.PreSaleOrderInfo
	var orderUser pty.PreSaleUserInfo
	var orderMarket pty.PreSaleMarketInfo
	keyOrderId := proKeyOrderId(orderId)
	orderInfo, err := a.db.Get(keyOrderId)
	if err != nil {
		clog.Error("presale userRecipt", "hash", a.txhash, "get kvdb failed",
			keyOrderId, "err", err)
		return nil, err
	}
	err1 := types.Decode(orderInfo, &order)
	if err1 != nil {
		clog.Error("presale userRecipt", "hash", a.txhash, "decode order failed",
			orderInfo, "err", err1)
		return nil, err1
	}
	if order.MarketStatus != "已发货" {
		return nil, pty.ErrorOrder
	}
	if order.Uid != uid {
		return nil, pty.ErrorOprea
	}
	marketId := order.MarketId
	keyMarket := proKeyMarketId(marketId)
	keyUser := proKeyUserId(uid)
	marketInfo, err3 := a.db.Get(keyMarket)
	if err3 != nil {
		return nil, err3
	}
	err = types.Decode(marketInfo, &orderMarket)
	if err != nil {
		return nil, err
	}
	userInfo, err2 := a.db.Get(keyUser)
	if err2 != nil {
		return nil, err2
	}
	err = types.Decode(userInfo, &orderUser)
	if err != nil {
		return nil, err
	}

	for _, info := range orderMarket.Currency {
		if order.Book != nil && info.CurrencyName == order.Book.CurrencyName {
			info.Frozen -= order.Book.Amount
			info.UseAble += order.Book.Amount
			if order.LastPay != nil {
				info.Frozen -= order.LastPay.Amount
				info.UseAble += order.LastPay.Amount
			}
		}
	}
	if orderMarket.IsTokenSup {
		orderMarket.Token.UseAble += order.TokenAmount
		for _, res := range orderUser.Token {
			if res.MarketId == marketId {
				res.TokenFrozen -= order.TokenAmount
				if res.TokenFrozen < 0 {
					return nil, pty.ErrorUserFrozen
				}
			}
		}
	}

	order.UserStatus = "已收货"
	order.OperateTime = time
	order.Status = "已收货"
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	valueMarket := types.Encode(&orderMarket)
	valueUser := types.Encode(&orderUser)
	//上链
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyOrderId, Value:value},
		&types.KeyValue{Key:keyUser, Value:valueUser}, &types.KeyValue{Key:keyMarket, Value:valueMarket})
	log := &pty.OrderInfoPreSale{OrderInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserRecipt, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execOverDue(overDue *pty.PreSaleOverDue, pubkey []byte) (*types.Receipt, error) {
	orderIdSlice := overDue.OrderId
	time := overDue.Time
	//time:=overDue.Time

	// err := a.JudgeUid(uid, pubkey)
	// if err != nil {
	// 	return nil, err
	// }
	adminpubkey, _ := hex.DecodeString(adminPubkey)
	if !bytes.Equal(adminpubkey, pubkey) {
		clog.Error("presale overdue", "hash", a.txhash, "adminpubkey:=\n",
			adminpubkey, "pubkey:=", pubkey)
		return nil, pty.ErrorOprea
	}
	//marketid和marketname是否一致 todo
	//两种价格上建立绑定
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	set := make(map[string]interface{})
	for _, value := range orderIdSlice {
		orderId := value
		// var ErrOrderId=errors.New(orderId+"")
		var order pty.PreSaleOrderInfo
		keyOrderId := proKeyOrderId(orderId)
		id := string(orderId)
		// keyUser := proKeyUserId(uid)
		orderInfo, err1 := a.db.Get(keyOrderId)
		if err1 != nil {
			clog.Error("presale overDue", "hash", a.txhash, "get kv failed", err1)
			err := errors.Wrap(pty.ErrorOrderNotExist, id)
			return nil, err
		}
		err2 := types.Decode(orderInfo, &order)
		if err2 != nil {
			clog.Error("presale overDue", "hash", a.txhash, "encode failed ", err2)
			err := errors.Wrap(pty.ErrorOrderNotExist, id)
			return nil, err
		}
		// if order.UserStatus != "已预订" {
		// 	return nil, pty.ErrorOrderStatus
		// }
		order.UserStatus = "订单过期"
		order.Status = "用户添加订单过期"
		order.OperateTime = time

		keyMarket := proKeyMarketId(order.MarketId)
		var marketInfo *pty.PreSaleMarketInfo
		var err error
		if info, ok := set[string(keyMarket)]; !ok {
			marketInfo, err = a.getMarketInfo(order.MarketId)
			if err != nil {
				err := errors.Wrapf(pty.ErrorOrderNotExist, id)
				return nil, err
			}
		} else {
			infoMarket, _ := info.(*pty.PreSaleMarketInfo)
			marketInfo = infoMarket
		}
		// marketInfo.Eth.UseAble += order.EthAmount
		// marketInfo.Bty.UseAble += order.BtyAmount
		// marketInfo.Token.UseAble += order.TokenAmount
		// marketInfo.Eth.Frozen -= order.EthAmount
		// marketInfo.Bty.Frozen -= order.BtyAmount
		for _, info := range marketInfo.Currency {
			if info.CurrencyName == order.Book.CurrencyName {
				info.Frozen -= order.Book.Amount
				info.UseAble += order.Book.Amount
			}
		}
		if marketInfo.IsTokenSup {
			marketInfo.Token.UseAble += order.TokenAmount
		}

		keyUser := proKeyUserId(order.Uid)
		var userInfo *pty.PreSaleUserInfo
		if info, ok := set[string(keyUser)]; !ok {
			userInfo, err = a.getUserInfo(order.Uid)
			if err != nil {
				return nil, err
			}
		} else {
			infoUser, _ := info.(*pty.PreSaleUserInfo)
			userInfo = infoUser
		}

		for _, res := range userInfo.Token {
			if res.MarketId == order.MarketId {
				res.TokenFrozen -= order.TokenAmount
				if res.TokenFrozen < 0 {
					err := errors.Wrapf(pty.ErrorUserFrozen, id)
					return nil, err
				}
			}
		}
		key := calcOrderKey(string(a.txhash))
		value := types.Encode(&order)
		// valueUser := types.Encode(userInfo)
		// valueMarket := types.Encode(marketInfo)
		set[string(keyMarket)] = marketInfo
		set[string(keyUser)] = userInfo
		kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyOrderId, Value:value})
	}
	var keys []string
	for k, _ := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch set[k].(type) {
		case *pty.PreSaleMarketInfo:
			v, _ := set[k].(*pty.PreSaleMarketInfo)
			kvv := types.Encode(v)
			kv = append(kv, &types.KeyValue{Key:[]byte(k), Value:kvv})
		case *pty.PreSaleUserInfo:
			v, _ := set[k].(*pty.PreSaleUserInfo)
			kvv := types.Encode(v)
			kv = append(kv, &types.KeyValue{Key:[]byte(k), Value:kvv})
		}
	}
	log := &pty.OverDuePreSale{Order: orderIdSlice}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleOverDue, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) getUserInfo(uid int64) (*pty.PreSaleUserInfo, error) {
	var userInfo pty.PreSaleUserInfo
	keyUser := proKeyUserId(uid)
	value, err := a.db.Get(keyUser)
	if err != nil {
		return &userInfo, err
	}
	err = types.Decode(value, &userInfo)
	if err != nil {
		return &userInfo, err
	}
	return &userInfo, nil

}

func (a *Action) getMarketInfo(marketId int64) (*pty.PreSaleMarketInfo, error) {
	var marketInfo pty.PreSaleMarketInfo
	keyMarket := proKeyMarketId(marketId)
	value, err := a.db.Get(keyMarket)
	if err != nil {
		return &marketInfo, err
	}
	err = types.Decode(value, &marketInfo)
	if err != nil {
		return &marketInfo, err
	}
	return &marketInfo, nil
}

// func (a *Action) getOrder(info *pty.PreSaleUserRec) (*pty.PreSaleUserRec, error) {
// 	b := 0
// 	for _, token := range info.Token {
// 		switch info.Currency {
// 		case "bty":
// 			token.UseAble_BTY += info.Amount
// 			b++
// 		case "eth":
// 			token.UseAble_ETH += info.Amount
// 			b++
// 		default:
// 			return nil, pty.ErrorToken
// 		}
// 	}
// 	if b == 0 {
// 		tokenInfo := &pty.PreSaleToken{}
// 		tokenInfo.MarketId = info.MarketId
// 		switch info.Currency {
// 		case "bty":
// 			tokenInfo.UseAble_BTY += info.Amount
// 		case "eth":
// 			tokenInfo.UseAble_ETH += info.Amount
// 		default:
// 			return nil, pty.ErrorToken
// 		}
// 		info.Token = append(info.Token, tokenInfo)
// 	}
// 	return &info, nil
// }

func (a *Action) execWithDraw(withDraw *pty.PreSaleWithDraw, pubkey []byte) (*types.Receipt, error) {
	isMarket := withDraw.IsMarket
	if isMarket {
		reply, err := a.execMarketDraw(withDraw, pubkey)
		if err != nil {
			return nil, err
		}
		return reply, nil
	} else {
		reply, err := a.execUserDraw(withDraw, pubkey)
		if err != nil {
			return nil, err
		}
		return reply, nil
	}
}

func (a *Action) isTrue(uid int64, pubkey []byte, message string) ([]byte, error) {
	keyJudge := proKeyUserId(uid)
	kvUid, err := a.db.Get(pubkey)
	if err != nil {
		clog.Error(message, "hash", a.txhash, "get uid failed", pubkey)
		return nil, pty.ErrorMatchMessage
	}
	if bytes.Equal(kvUid, keyJudge) {
		clog.Error(message, "hash", a.txhash, "id not match to pubkey", kvUid)
		return nil, pty.ErrorMatchMessage
	}
	return kvUid, nil
}

func (a *Action) JudgeUid(uid int64, pubkey []byte) error {
	keyJudge := proKeyUidJudge(uid)
	value, err := a.db.Get(keyJudge)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, value) {
		return pty.ErrorMatchMessage
	}
	return nil
}

func (a *Action) JudgeMarkeyId(marketId int64, pubkey []byte) error {
	keyJudge := proKeyMarketIdJudge(marketId)
	value, err := a.db.Get(keyJudge)
	if err != nil {
		return err
	}
	if !bytes.Equal(pubkey, value) {
		return pty.ErrorMatchMessage
	}
	return nil
}

func (a *Action) execMarketDraw(withDraw *pty.PreSaleWithDraw, pubkey []byte) (*types.Receipt, error) {
	marketid := withDraw.Id
	amount := withDraw.Amount
	toAddress := withDraw.ToAddress
	currency := withDraw.Currency
	isToken := withDraw.IsToken
	time := withDraw.Time
	minerFee := withDraw.MinerFee
	err := a.JudgeMarkeyId(marketid, pubkey)
	if err != nil {
		return nil, err
	}
	keyMarket := proKeyMarketId(marketid)
	marketInfo, err := a.getMarketInfo(marketid)
	if err != nil {
		return nil, err
	}
	// if toAddress != marketInfo.Eth.Address && toAddress != marketInfo.Bty.Address {
	// 	return nil, pty.ErrorAddress
	// }

	if isToken {
		marketInfo.Token.UseAble -= amount
		if marketInfo.Token.UseAble < 0 {
			return nil, pty.ErrorAmount
		}
		opera := &pty.PreSaleOpreation{}

		opera.Time = time
		opera.Status = "用户提Token"
		opera.ToAddress = toAddress
		opera.MinerFee = minerFee

		marketInfo.Operate = opera
	} else {
		for _, info := range marketInfo.Currency {
			if info.CurrencyName == currency {
				info.UseAble -= amount
			}
			// if strings.Compare(info.Address, toAddress) != 0 {
			// 	return nil, pty.ErrorAddress
			// }
			opera := &pty.PreSaleOpreation{}

			opera.Time = time
			opera.Status = "用户提币"
			opera.ToAddress = toAddress
			opera.MinerFee = minerFee

			marketInfo.Operate = opera
			if info.UseAble < 0 {
				return nil, pty.ErrorAmount
			}
		}
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(marketInfo)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyMarket, Value:value})
	log := &pty.WithDrawPS{WithDrawInfo: withDraw}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleWithDraw, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserDraw(withDraw *pty.PreSaleWithDraw, pubkey []byte) (*types.Receipt, error) {
	uid := withDraw.Id
	amount := withDraw.Amount
	toAddress := withDraw.ToAddress
	currency := withDraw.Currency
	time := withDraw.Time
	minerFee := withDraw.MinerFee
	err := a.JudgeUid(uid, pubkey)
	if err != nil {
		return nil, err
	}
	keyUser := proKeyUserId(uid)
	userInfo, err := a.getUserInfo(uid)
	if err != nil {
		return nil, err
	}
	// if toAddress != userInfo.Eth.Address && toAddress != userInfo.Bty.Address {
	// 	return nil, pty.ErrorAddress
	// }
	// switch currency {
	// case "bty":
	// 	userInfo.Bty.UseAble -= amount
	// 	userInfo.Operate.Time = time
	// 	userInfo.Operate.ToAddress = toAddress
	// 	userInfo.Operate.Status = "用户提BTY"
	// case "eth":
	// 	userInfo.Eth.UseAble -= amount
	// 	userInfo.Operate.Time = time
	// 	userInfo.Operate.ToAddress = toAddress
	// 	userInfo.Operate.Status = "用户提ETH"
	// default:
	// 	return nil, pty.ErrorToken
	// }
	for _, info := range userInfo.Currency {
		if info.CurrencyName == currency {
			info.UseAble -= amount
		}
		// if strings.Compare(info.Address, toAddress) != 0 {
		// 	return nil, pty.ErrorAddress
		// }
		if info.UseAble < 0 {
			return nil, pty.ErrorAmount
		}
	}
	opera := &pty.PreSaleOpreation{}

	opera.Time = time
	opera.Status = "用户提币"
	opera.ToAddress = toAddress
	opera.MinerFee = minerFee

	userInfo.Operate = opera
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(userInfo)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyUser, Value:value})
	log := &pty.WithDrawPS{WithDrawInfo: withDraw}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleWithDraw, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execChange(changePd *pty.PreSaleChangePd, pubkey []byte) (*types.Receipt, error) {
	// pdName := changePd.PdName
	// pdId := changePd.PdId
	// price := changePd.Price
	// deposit := changePd.Deposit
	// psEnd := changePd.PsEnd
	// psStart := changePd.PsStart
	// puStart := changePd.PuStart
	// puEnd := changePd.PuEnd
	// marketId := changePd.MarketId
	// marketName := changePd.MarketName
	adminId := changePd.AdminId
	// status := changePd.Status
	isDelete := changePd.IsDelete

	// keyAdmin:=proKeyUidJudge(adminId)
	err := a.JudgeMarkeyId(adminId, pubkey)
	if err != nil {
		return nil, err
	}
	if isDelete {
		return a.deleteDatabase(changePd)
	} else {
		return a.changeDatabase(changePd)
	}

}

func (a *Action) execUserTokenWithDraw(withDraw *pty.PreSaleUserTokenWithDraw, pubkey []byte) (*types.Receipt, error) {

	err := a.JudgeUid(withDraw.Uid, pubkey)
	if err != nil {
		return nil, err
	}
	keyUser := proKeyUserId(withDraw.Uid)
	userInfo, err := a.getUserInfo(withDraw.Uid)
	if err != nil {
		return nil, err
	}

	for _, info := range userInfo.Token {
		if info.MarketId == withDraw.MarketId {
			info.TokenFrozen -= withDraw.Amount
		}

		if info.TokenFrozen < 0 {
			return nil, pty.ErrorAmount
		}
	}
	opera := &pty.PreSaleOpreation{}

	opera.Time = 0
	opera.Status = "用户提币"
	opera.ToAddress = withDraw.Address

	userInfo.Operate = opera
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(userInfo)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyUser, Value:value})
	log := &pty.UserInfoSale{UserInfo:userInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserTokenWithDraw, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) execUserShip(userShip *pty.PreSaleUserShip, pubkey []byte) (*types.Receipt, error) {

	err := a.JudgeUid(userShip.Uid,pubkey)
	if err != nil {
		return nil, err
	}
	userInfo ,err := a.getUserInfo(userShip.Uid)
	if err != nil {
		return nil, err
	}

	var order pty.PreSaleOrderInfo
	keyOrderId := proKeyOrderId(userShip.Oid)
	orderInfo, err := a.db.Get(keyOrderId)
	if err != nil {
		clog.Error("presale userShip", "hash", a.txhash, "get kvdb failed",
			keyOrderId, "err", err)
		return nil,err
	}
	err1 := types.Decode(orderInfo, &order)
	if err1 != nil {
		clog.Error("presale userShip", "hash", a.txhash, "decode order failed",
			orderInfo, "err", err)
		return nil, err
	}

	order.MarketStatus = "待发货"
	order.UserStatus = "用户提货"
	order.Status = "用户提货"
	c := 0
	for _,addr := range userInfo.Address {
		if addr.AddressId == userShip.AddressId {
			order.Address = addr.Address
			c++
		}
	}
	if c == 0 {
		return nil,errors.New("address not exists")
	}
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(&order)
	//上链
	a.db.Set(keyOrderId, value)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyOrderId, Value:value})
	log := &pty.OrderInfoPreSale{OrderInfo: &order}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleUserShip, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (a *Action) deleteDatabase(changePd *pty.PreSaleChangePd) (*types.Receipt, error) {
	pdInfo, err := a.getPdInfo(changePd.PdId)
	if err != nil {
		return nil, err
	}
	pdInfo.Status = "已删除"

	keyPd := proKeyPdId(changePd.PdId)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(pdInfo)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyPd, Value:value})
	log := &pty.PdInfoPresale{PdInfo: pdInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleChangePd, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil

}

func (a *Action) getPdInfo(pdId int64) (*pty.PresalePdInfo, error) {
	key := proKeyPdId(pdId)
	value, err := a.db.Get(key)
	var pdInfo pty.PresalePdInfo
	if err != nil {
		return &pdInfo, pty.ErrorPdNotExist
	}
	err = types.Decode(value, &pdInfo)
	if err != nil {
		return &pdInfo, err
	}
	return &pdInfo, nil
}

func (a *Action) changeDatabase(changePd *pty.PreSaleChangePd) (*types.Receipt, error) {
	pdInfo, err := a.getPdInfo(changePd.PdId)
	if err != nil {
		return nil, err
	}
	pdInfo.PdName = changePd.PdName
	pdInfo.PdId = changePd.PdId
	pdInfo.Price = changePd.Price
	pdInfo.Deposit = changePd.Deposit
	pdInfo.PsStart = changePd.PsStart
	pdInfo.PsEnd = changePd.PsEnd
	pdInfo.PuStart = changePd.PuStart
	pdInfo.PuEnd = changePd.PuEnd
	pdInfo.MarketId = changePd.MarketId
	pdInfo.MarketName = changePd.MarketName
	pdInfo.AdminId = changePd.AdminId
	pdInfo.Status = changePd.Status

	keyPd := proKeyPdId(changePd.PdId)
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	key := calcOrderKey(string(a.txhash))
	value := types.Encode(pdInfo)
	kv = append(kv, &types.KeyValue{Key:[]byte(key), Value:value}, &types.KeyValue{Key:keyPd, Value:value})
	log := &pty.PdInfoPresale{PdInfo: pdInfo}
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyLogPreSaleChangePd, Log: types.Encode(log)})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}
