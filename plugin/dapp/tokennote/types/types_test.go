// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/hex"
	"testing"

	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common"
	"github.com/stretchr/testify/assert"

	rpctypes "github.com/33cn/chain33/rpc/types"
	"fmt"
	"gitlab.33.cn/chain33/chain33/common/address"
	"github.com/decred/base58"
	"github.com/33cn/chain33/common/crypto"
)


func TestToBytes(t *testing.T) {
	var req TokennoteAction
	req.Ty = TokennoteActionCreate
	req.Value = &TokennoteAction_TokennoteCreate{
		TokennoteCreate:&TokennoteCreate{
			Issuer:"issuer",
			IssuerName:"杨成涛",
			IssuerPhone:"13716976706",
			IssuerId:"1111111",
			Acceptor:"acceptor",
			Balance:6,
			Currency:"7",
			AcceptanceDate:1550028373,
			Rate:1200,
			Introduction:"this is a introduction",
		},
	}

	var tx types.Transaction
	var priv crypto.PrivKey
	priv = getprivkey("72c3879f1f9b523f266a9545b69bd41c0251483a93e21e348e85118afe17a5e2")
	tx.Execer = []byte("tokennote")
	tx.To = "1McCR6oHLg5KLbgbtPCvvsn8uYNvEZsQd4"
	tx.Expire = 0
	tx.Fee = 1e6
	tx.Nonce = 123456789
	tx.Payload = types.Encode(&req)

	tx.Sign(types.SECP256K1,priv)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(&tx)))
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New("secp256k1")
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func TestBase58(t *testing.T) {
	addr := "15M1iwuk5eGbqngsjUbzNosTjdiVq7ryHc"
	byteaddr := base58.Decode(addr)
	fmt.Println("bytes:",hex.EncodeToString(byteaddr))
	addr1 := base58.Encode(byteaddr)
	fmt.Println("addr1:",addr1)
}

func TestExecName(t *testing.T) {
	fmt.Println("addr:",address.ExecAddress("token")," para:",address.ExecAddress("user.p.developer.token"))
	pub,_ := common.FromHex("0307b368e8c2e9a7f9ed5f4b85a845ee6cae9fac7b9743e1dddfa823161284a13c")
	p := make([]byte,len(pub))
	fmt.Println(p)
	copy(p,pub[:])
	fmt.Println(p)
	add := address.PubKeyToAddress(pub).String()
	fmt.Println("addr:",add," version:",address.PubKeyToAddress(pub).Version)
	fmt.Println("pub:",pub)
}

func TestPayload(t *testing.T) {
	tx := "0a09746f6b656e6e6f7465125e0a590a066973737565721209e69da8e68890e6b69b1a0b31333731363937363730362207313131313131312a086163636570746f7230063a013740d59c8ee30548b009521674686973206973206120696e74726f64756374696f6e9806071ab7030801122103bfb792eb89fe7103b599252896f6817657dcb50a37bd89865671bac8f42cae981a8f03d1ad3def8e9fe9beb9e9ee9ee9fef8eb9d76e5ed1ae7dd1ad3aebdef7ef7ef9eb9ef6d76d3d7baf5d6bc7baf3cf747ba6faf5bd5ad1bdf5df7dfbdf5dfadfddfbdfadfbdf4dfadb6d3bdf5df5df5df5df5df5df5d9ad3ceb5eb7eb7eb9ef4ef8e9fef6df4d3addad35dfbe34779f5cf1e7b7d39e3c6f4d3de76d7aef8ebcebdef7db4ebdef7db4eb5db4ebde9eef8ef6e9feb8ef9eb7ef8ebde9fe9ef7cd3ad3bd5ae9dd3cd35d76db5d376df6fbf7679bf3d7deef5d376f9f7ddb9dbcf7a7faf35efae7b75c6f9d1adfb6ddf3df3ae7aef56da73c7f8d9c69ef7cd5ae3adf4e38d36db4e77eda6fbd5b6dce7ddba7db6b5f3adfd6fcd5c73ae9b75d779e5c6fce35ef769fd5d71fdf67b9f39e9e6db6b8f5ee75d36db4d1af757b8e3cf1f6fbd1bd75778eb873ae7569d75ed9ddb7e9ad37f1a79dd9fdfc7fa75ef79735e9b7f9db7dba79fdbddbc6df7daf787b7d39df4f3d6fa77a6f4d37ddadb6df5e1deb7e37e76dfae9fe3ce1cebbdf9e1be1ceb6ebbeb6ef8e74e37efaefaef7e9edfcef9e7de1eefae39e5aef7e75eb8df83089b6d6b0033a22314d634352366f484c67354b4c6267627450437676736e3875594e76455a73516434"
	bt ,err := common.FromHex(tx)
	if err != nil {
		panic(err)
	}
	var tt types.Transaction
	err = types.Decode(bt,&tt)
	if err != nil {
		panic(err)
	}
	fmt.Println("pubkey",hex.EncodeToString(tt.Signature.Pubkey))
	fmt.Println(" addr: ",address.PubKeyToAddress(tt.Signature.Pubkey))
	fmt.Println("detail",string(tt.Execer))
	//var payload types2.TraceplatformAction
	//err = types.Decode(tt.Payload,&payload)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(payload.GetTraceplatformAddGood().Goodinfo.Goodinfo)

	//hex.EncodeToString(tt.Hash())
	//ttt ,_:= tt.GetTxGroup()
	//fmt.Println("ttt:",ttt)



	//var  new types.Transaction
	//types.Decode(tt.Header,&new)
	//
	//var newp TokenAction
	//err = types.Decode(new.Payload,&newp)
	//if err != nil {
	//	fmt.Println("EWR:",err)
	//}
	//fmt.Println("new :",newp)
}

func isExecAddrMatch(name string, to string) bool {
	toaddr := address.ExecAddress(name)
	return toaddr == to
}

func TestDecodeLogTokenTransfer(t *testing.T) {
	var logTmp = &types.ReceiptAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteTransfer,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenTransfer", result.Logs[0].TyName)
}

func TestDecodeLogTokenDeposit(t *testing.T) {
	var logTmp = &types.ReceiptAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteDeposit,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenDeposit", result.Logs[0].TyName)
}

func TestDecodeLogTokenExecTransfer(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteExecTransfer,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenExecTransfer", result.Logs[0].TyName)
}

func TestDecodeLogTokenExecWithdraw(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteExecWithdraw,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenExecWithdraw", result.Logs[0].TyName)
}

func TestDecodeLogTokenExecDeposit(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteExecDeposit,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenExecDeposit", result.Logs[0].TyName)
}

func TestDecodeLogTokenExecFrozen(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteExecFrozen,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   5,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenExecFrozen", result.Logs[0].TyName)
}

func TestDecodeLogTokenExecActive(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteExecActive,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   0,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenExecActive", result.Logs[0].TyName)
}

func TestDecodeLogTokenGenesisTransfer(t *testing.T) {
	var logTmp = &types.ReceiptAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteGenesisTransfer,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   1,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenGenesisTransfer", result.Logs[0].TyName)
}

func TestDecodeLogTokenGenesisDeposit(t *testing.T) {
	var logTmp = &types.ReceiptExecAccountTransfer{}
	dec := types.Encode(logTmp)
	strdec := hex.EncodeToString(dec)
	rlog := &rpctypes.ReceiptLog{
		Ty:  TyLogTokennoteGenesisDeposit,
		Log: "0x" + strdec,
	}

	logs := []*rpctypes.ReceiptLog{}
	logs = append(logs, rlog)

	var data = &rpctypes.ReceiptData{
		Ty:   2,
		Logs: logs,
	}
	result, err := rpctypes.DecodeLog([]byte("token"), data)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "LogTokenGenesisDeposit", result.Logs[0].TyName)
}
