package executor

import (
	pty "github.com/33cn/plugin/plugin/dapp/presale/types"
	"github.com/33cn/chain33/types"
)

func (p *PreSale) Query_GetPreSaleUserInfo(in *types.Int64) (types.Message,error) {
	reply, err := p.getUserInfo(in.Data)
	if err != nil {
		return nil, err
	}
	return reply, nil

}

func (p *PreSale) Query_GetPreSaleMarketInfo(in *types.Int64) (types.Message,error) {
	reply, err := p.getMarketInfo(in.Data)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (p *PreSale) Query_GetGoldbillPlatform(in *types.Int64) (types.Message,error) {
	reply, err := p.getOrderInfo(in.Data)
	if err != nil {
		return nil, err
	}
	return reply, nil

}

func (p *PreSale) Query_GetPreSalePdInfo(in *types.Int64) (types.Message,error) {
	reply, err := p.getPdInfo(in.Data)
	if err != nil {
		return nil, err
	}
	return reply, nil

}

func (p *PreSale) getUserInfo(uid int64) (*pty.UserInfoSale, error) {
	keyGet := proKeyUserId(uid)
	value, err := p.GetStateDB().Get(keyGet)
	if err != nil {
		clog.Info("GetPreSaleUserInfo", "get db key failed", "not found")
		return nil, err
	}
	var userInfo pty.PreSaleUserInfo
	var reply pty.UserInfoSale
	err = types.Decode(value, &userInfo)
	if err != nil {
		clog.Info("GetPreSaleUserInfo", "decode failed", "not found")
		return nil, err
	}
	reply.UserInfo = &userInfo
	return &reply, nil
}

func (p *PreSale) getMarketInfo(marketId int64) (*pty.MarketInfoPreSale, error) {
	keyGet := proKeyMarketId(marketId)
	value, err := p.GetStateDB().Get(keyGet)
	if err != nil {
		clog.Info("GetPreSaleMarketInfo", "get db key failed", "not found")
		return nil, err
	}
	var marketInfo pty.PreSaleMarketInfo
	var reply pty.MarketInfoPreSale
	err = types.Decode(value, &marketInfo)
	if err != nil {
		clog.Info("GetPreSaleMarketInfo", "decode failed", "not found")
		return nil, err
	}
	reply.MarketInfo = &marketInfo
	return &reply, nil
}

func (p *PreSale) getOrderInfo(orderId int64) (*pty.OrderInfoPreSale, error) {
	keyGet := proKeyOrderId(orderId)
	value, err := p.GetStateDB().Get(keyGet)
	if err != nil {
		clog.Info("GetPreSaleOrderInfo", "get db key failed", "not found")
		return nil, err
	}
	var orderInfo pty.PreSaleOrderInfo
	var reply pty.OrderInfoPreSale
	err = types.Decode(value, &orderInfo)
	if err != nil {
		clog.Info("GetPreSaleOrderInfo", "decode failed", "not found")
		return nil, err
	}
	reply.OrderInfo = &orderInfo
	return &reply, nil
}

func (p *PreSale) getPdInfo(pdId int64) (*pty.PdInfoPresale, error) {
	keyGet := proKeyPdId(pdId)
	value, err := p.GetStateDB().Get(keyGet)
	if err != nil {
		clog.Info("GetPreSalePdInfo", "get db key failed", "not found")
		return nil, err
	}
	var pdInfo pty.PresalePdInfo
	var reply pty.PdInfoPresale
	err = types.Decode(value, &pdInfo)
	if err != nil {
		clog.Info("GetPreSalePdInfo", "decode failed", "not found")
		return nil, err
	}
	reply.PdInfo = &pdInfo
	return &reply, nil
}