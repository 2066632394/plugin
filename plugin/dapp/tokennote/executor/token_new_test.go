// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	cty "github.com/33cn/chain33/system/dapp/coins/types"
	"github.com/33cn/chain33/types"
	tokenty "github.com/33cn/plugin/plugin/dapp/token/types"
	tokennotety "github.com/33cn/plugin/plugin/dapp/tokennote/types"
	"google.golang.org/grpc"
	"encoding/hex"
	"math/big"
	"math"
)

var (
	isMainNetTest = true
	isParaNetTest bool
)

var (
	mainNetgrpcAddr = "172.16.100.202:8802"
	mainNetgrpcAddr1 = "jiedian1.bityuan.com:8802"
	ParaNetgrpcAddr = "172.16.100.202:8802"
	mainClient      types.Chain33Client
	paraClient      types.Chain33Client
	//jsonRPCClient   *rpc.JSONClient
	r               *rand.Rand

	ErrTest = errors.New("ErrTest")

	addrexec     string
	addr         string
	privkey      = getprivkey("0x016959665e2dd245ffef62a73df1dfe633308e57421366e260da136731137961")
	lendPrivkey  = getprivkey("0x81b23fad7975146167bd9e98ca4e925707bfbaac21dd652e4e5c3aaec9062867")
	to1Privkey   = getprivkey("0x965fa3de8b4c85d6f92d850b0984fa50f001e2bd6a58d8b7c55ab2c86292d17c")
	ccnyPrivkey  = getprivkey("0xc20e3eca2d63855cc3d24721b02bd937f6dd1cb56b62e599c507f317c0f6e950")
	szhPrivkey  = getprivkey("0x5ced1e2528e6063d0a447397c11f6e92c937d9a9771a083a67b80b53ea976c64")
	privGenesis  crypto.PrivKey
	privkeySuper crypto.PrivKey
)

const (
	//defaultAmount = 1e10
	fee = 1e6
)

//for token
var (
	tokenName   = "NEW"
	tokenSym    = "NEW"
	tokenIntro  = "newtoken"
	tokenPrice  int64
	tokenAmount int64 = 1000 * 1e4 * 1e4
	//execName          = "user.p.loantest.tokennote"
	execName          = "tokennote"
	feeForToken int64 = 1e5
	transToAddr       = "1NYxhca2zVMzxFqMRJdMcZfrSFnqbqotKe"
	transAmount int64 = 100 * 1e4 * 1e4
	walletPass        = "123456"
	issuer = "1F2MJDXevhAjsfZKTNPwkwdrDzVvW5a8T4"
	issuerName = "szh"
	issuerPhone = "13754338419"
	issuerId = "33333333333333333"
	to 	   = "1NSXuM85ochWHYUVJjmKeWi4ChrMUCwHFr"
	to1    = "1JBpM7TWLFCGQmUF9rqsgii1mcFRzcbEcV"
	ccnyIssuer = "16dzhkkZ2dtRysMEarGKkSGWPCssEZku9c"
	szh = "19ndKJEKLEcXzg8nEpkKbwJo6566WUuhtf"


)

//测试过程：
//1. 初始化账户，导入有钱的私钥，创建一个新账户，往这个新账户打钱（用来签名和扣手续费）
//2. 产生precreate的一种token
//3. finish这个token
//4. 向一个地址转账token
//5. 可选：在平行链上进行query

func init() {
	fmt.Println("Init start")
	defer fmt.Println("Init end")

	//if !isMainNetTest && !isParaNetTest {
	//	return
	//}
	maxReceLimit := grpc.WithMaxMsgSize(30*1024*1024)
	conn, err := grpc.Dial(mainNetgrpcAddr, grpc.WithInsecure(),maxReceLimit)
	if err != nil {
		panic(err)
	}
	mainClient = types.NewChain33Client(conn)

	conn, err = grpc.Dial(ParaNetgrpcAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	paraClient = types.NewChain33Client(conn)

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	addrexec = address.ExecAddress("user.p.loantest.tokennote")
	//jsonRPCClient ,err = rpc.NewJSONClient("http://172.16.100.202:8801")
	//if err != nil {
	//	panic(err)
	//}
	privGenesis = getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944")
	privkeySuper = getprivkey("4a92f3700920dc422c8ba993020d26b54711ef9b3d74deab7c3df055218ded42")
}

func TestCheckAddr(t *testing.T) {
	err := address.CheckAddress(to)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("ok")
}

func TestInitAccount(t *testing.T) {
	if !isMainNetTest {
		fmt.Println("fail",isMainNetTest)
		return
	}
	fmt.Println("TestInitAccount start")
	defer fmt.Println("TestInitAccount end")

	//need update to fixed addr here
	//addr = ""
	//privkey = ""
	//addr, privkey = genaddress()
	label := strconv.Itoa(int(types.Now().UnixNano()))
	params := types.ReqWalletImportPrivkey{Privkey: common.ToHex(privkey.Bytes()), Label: label}

	unlock := types.WalletUnLock{Passwd: walletPass, Timeout: 0, WalletOrTicket: false}
	_, err := mainClient.UnLock(context.Background(), &unlock)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)

	_, err = mainClient.ImportPrivkey(context.Background(), &params)
	if err != nil && err != types.ErrPrivkeyExist {
		fmt.Println(err)
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)
	/*
		txhash, err := sendtoaddress(mainClient, privGenesis, addr, defaultAmount)

		if err != nil {
			t.Error(err)
			return
		}
		if !waitTx(txhash) {
			t.Error(ErrTest)
			return
		}

		time.Sleep(5 * time.Second)
	*/
}

func TestCreate(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestPrecreate start")
	defer fmt.Println("TestPrecreate end")
	v := &tokennotety.TokennoteCreate{
		Issuer:       szh,
		IssuerName:issuerName,
		IssuerId:issuerId,
		IssuerPhone:issuerPhone,
		Currency:"SZH000014",
		//Acceptor:    issuer,
		AcceptanceDate:1554912000,
		Introduction: tokenIntro,
		Balance:        10000000,
		Rate:        28,
		Repayamount:        0,
	}
	create := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionCreate,
		Value: &tokennotety.TokennoteAction_TokennoteCreate{TokennoteCreate: v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(create),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, szhPrivkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)

		t.Error(err)
		return
	}
	fmt.Println("replay:",reply,common.ToHex(reply.Msg))

	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}
	fmt.Println("txhash:",hex.EncodeToString(tx.Hash()))
	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}
	time.Sleep(5 * time.Second)

}

func TestLoan(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestLoan start")
	defer fmt.Println("TestLoan end")

	v := &tokennotety.TokennoteLoan{Symbol: "SZH0004", To: szh,Amount:10000000}
	finish := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionLoan,
		Value: &tokennotety.TokennoteAction_TokennoteLoan{TokennoteLoan: v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(finish),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, privkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
	//reply, err := mainClient.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	fmt.Println("err", err)
	//	t.Error(err)
	//	return
	//}
	//if !reply.IsOk {
	//	fmt.Println("err = ", reply.GetMsg())
	//	t.Error(ErrTest)
	//	return
	//}
	//
	//if !waitTx(tx.Hash()) {
	//	t.Error(ErrTest)
	//	return
	//}
	//time.Sleep(5 * time.Second)

}

func TestTransferToTokennote(t *testing.T) {
	fmt.Println("transfer to tokennote start")
	defer  fmt.Println("transfer to tokennote end ")
	v := &types.AssetsTransferToExec{Cointoken:"CNYY",Amount:1000000000,Note:[]byte("szh transfer cnyy to tokennote"),ExecName:"tokennote",To:to}
	transfer := &tokenty.TokenAction{
		Ty:tokenty.TokenActionTransferToExec,
		Value:&tokenty.TokenAction_TransferToExec{TransferToExec:v},
	}
	tx := &types.Transaction{
		Execer:  []byte("token"),
		Payload: types.Encode(transfer),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("token"),
	}
	tx.Sign(types.SECP256K1, lendPrivkey)
	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}
	time.Sleep(5 * time.Second)

}

func TestTransferToToken(t *testing.T) {
	fmt.Println("transfer to tokennote start")
	defer  fmt.Println("transfer to tokennote end ")
	v := &types.AssetsTransfer{Cointoken:"CCNY",Amount:100000000,Note:[]byte("szh transfer cnyy to tokennote"),To:"16dzhkkZ2dtRysMEarGKkSGWPCssEZku9c"}
	transfer := &tokennotety.TokennoteAction{
		Ty:tokennotety.ActionTransfer,
		Value:&tokennotety.TokennoteAction_Transfer{Transfer:v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(transfer),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, szhPrivkey)
	fmt.Println("realaddr : ",tx.GetRealToAddr())
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
	var reply *types.Reply
	//err := jsonRPCClient.Call("SendTransaction",tx,reply)
	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}
	time.Sleep(5 * time.Second)

}

type serverResponse struct {
	ID     uint64      `json:"id"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

func TestTransferWithdraw(t *testing.T) {
	fmt.Println("TransferWithdraw to tokennote start")
	defer  fmt.Println("TransferWithdraw to tokennote end ")
	v := &types.AssetsWithdraw{Cointoken:"CNYY",Amount:1000000000,Note:[]byte("szh transfer cnyy to tokennote"),To:"1McCR6oHLg5KLbgbtPCvvsn8uYNvEZsQd4"}
	transfer := &tokenty.TokenAction{
		Ty:tokenty.ActionWithdraw,
		Value:&tokenty.TokenAction_Withdraw{Withdraw:v},
	}
	tx := &types.Transaction{
		Execer:  []byte("token"),
		Payload: types.Encode(transfer),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("token"),
	}
	tx.Sign(types.SECP256K1, lendPrivkey)
	fmt.Println("realaddr : ",tx.GetRealToAddr())
	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}
	time.Sleep(5 * time.Second)

}

func TestLoanedAgree(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestLoanAgree start")
	defer fmt.Println("TestAgree end")

	v := &tokennotety.TokennoteLoanedAgree{Symbol: "SZH0003", Owner:issuer,Loantime:0}
	finish := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionLoanedAgree,
		Value: &tokennotety.TokennoteAction_TokennoteLoanedAgree{TokennoteLoanedAgree: v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(finish),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, szhPrivkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
	//reply, err := mainClient.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	fmt.Println("err", err)
	//	t.Error(err)
	//	return
	//}
	//if !reply.IsOk {
	//	fmt.Println("err = ", reply.GetMsg())
	//	t.Error(ErrTest)
	//	return
	//}
	//
	//if !waitTx(tx.Hash()) {
	//	t.Error(ErrTest)
	//	return
	//}
	//time.Sleep(5 * time.Second)

}

func TestLoanedReject(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestLoanedReject start")
	defer fmt.Println("TestLoanedReject end")

	v := &tokennotety.TokennoteLoanedReject{Symbol: "SZH0004", Owner:issuer,LoanTime:1550215581}
	finish := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionLoanedReject,
		Value: &tokennotety.TokennoteAction_TokennoteLoanedReject{TokennoteLoanedReject: v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(finish),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, szhPrivkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))

	//reply, err := mainClient.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	fmt.Println("err", err)
	//	t.Error(err)
	//	return
	//}
	//if !reply.IsOk {
	//	fmt.Println("err = ", reply.GetMsg())
	//	t.Error(ErrTest)
	//	return
	//}
	//
	//if !waitTx(tx.Hash()) {
	//	t.Error(ErrTest)
	//	return
	//}
	//time.Sleep(5 * time.Second)

}

func TestLoanCashed(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestLoanCashed start")
	defer fmt.Println("TestLoanCashed end")

	v := &tokennotety.TokennoteCashed{Symbol: "SZH0003", Cash:12}
	finish := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionCashed,
		Value: &tokennotety.TokennoteAction_TokennoteCashed{TokennoteCashed: v},
	}
	tx := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(finish),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	tx.Sign(types.SECP256K1, privkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
	//reply, err := mainClient.SendTransaction(context.Background(), tx)
	//if err != nil {
	//	fmt.Println("err", err)
	//	t.Error(err)
	//	return
	//}
	//if !reply.IsOk {
	//	fmt.Println("err = ", reply.GetMsg())
	//	t.Error(ErrTest)
	//	return
	//}
	//
	//if !waitTx(tx.Hash()) {
	//	t.Error(ErrTest)
	//	return
	//}
	//time.Sleep(5 * time.Second)

}

func TestTransferToken(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestTransferToken start")
	defer fmt.Println("TestTransferToken end")

	v := &tokenty.TokenAction_Transfer{Transfer: &types.AssetsTransfer{Cointoken: tokenSym, Amount: transAmount, Note: []byte(""), To: transToAddr}}
	transfer := &tokenty.TokenAction{Value: v, Ty: tokenty.ActionTransfer}

	tx := &types.Transaction{Execer: []byte(execName), Payload: types.Encode(transfer), Fee: fee, To: addrexec}
	tx.Nonce = r.Int63()
	tx.Sign(types.SECP256K1, privkey)

	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}

}

func TestMint(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestMint start")
	defer fmt.Println("TestMint end")

	v := &tokennotety.TokennoteAction_TokennoteMint{
		TokennoteMint: &tokennotety.TokennoteMint{
			Symbol: "SZH000001",
			Amount: transAmount,
	}}
	mint := &tokennotety.TokennoteAction{Value: v, Ty: tokennotety.TokennoteActionMint}

	tx := &types.Transaction{
		Execer:  []byte("tokennote"),
		Payload: types.Encode(mint),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("tokennote"),
	}
	tx.Nonce = r.Int63()
	tx.Sign(types.SECP256K1, privkey)

	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}

}


func TestBurn(t *testing.T) {
	if !isMainNetTest {
		return
	}
	fmt.Println("TestBurn start")
	defer fmt.Println("TestBurn end")

	v := &tokennotety.TokennoteAction_TokennoteBurn{
		TokennoteBurn: &tokennotety.TokennoteBurn{
			Symbol: "SZH000001",
			Amount: 1e7,
		}}
	mint := &tokennotety.TokennoteAction{Value: v, Ty: tokennotety.TokennoteActionBurn}

	tx := &types.Transaction{
		Execer:  []byte("tokennote"),
		Payload: types.Encode(mint),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("tokennote"),
	}
	tx.Nonce = r.Int63()
	tx.Sign(types.SECP256K1, privkey)

	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(ErrTest)
		return
	}

	if !waitTx(tx.Hash()) {
		t.Error(ErrTest)
		return
	}

}

func TestQueryAsset(t *testing.T) {
	if !isParaNetTest {
		return
	}
	fmt.Println("TestQueryAsset start")
	defer fmt.Println("TestQueryAsset end")

	var req types.ChainExecutor
	req.Driver = execName
	req.FuncName = "GetAccountTokenAssets"

	var reqAsset tokenty.ReqAccountTokenAssets
	reqAsset.Address = addr
	reqAsset.Execer = execName

	req.Param = types.Encode(&reqAsset)

	reply, err := paraClient.QueryChain(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("Query reply err")
		t.Error(ErrTest)
		return
	}
	var res tokenty.ReplyAccountTokenAssets
	err = types.Decode(reply.Msg, &res)
	if err != nil {
		t.Error(err)
		return
	}
	for _, ta := range res.TokenAssets {
		//balanceResult := strconv.FormatFloat(float64(ta.Account.Balance)/float64(types.TokenPrecision), 'f', 4, 64)
		//frozenResult := strconv.FormatFloat(float64(ta.Account.Frozen)/float64(types.TokenPrecision), 'f', 4, 64)
		fmt.Println(ta.Symbol)
		fmt.Println(ta.Account.Addr)
		fmt.Println(ta.Account.Currency)
		fmt.Println(ta.Account.Balance)
		fmt.Println(ta.Account.Frozen)

	}

}

//***************************************************
//**************common actions for Test**************
//***************************************************
func sendtoaddress(c types.Chain33Client, priv crypto.PrivKey, to string, amount int64) ([]byte, error) {
	v := &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: amount}}
	transfer := &cty.CoinsAction{Value: v, Ty: cty.CoinsActionTransfer}
	tx := &types.Transaction{Execer: []byte("coins"), Payload: types.Encode(transfer), Fee: fee, To: to}
	tx.Nonce = r.Int63()
	tx.Sign(types.SECP256K1, priv)
	// Contact the server and print out its response.
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		return nil, err
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		return nil, errors.New(string(reply.GetMsg()))
	}
	return tx.Hash(), nil
}

func waitTx(hash []byte) bool {

	i := 0
	for {
		i++
		if i%100 == 0 {
			fmt.Println("wait transaction timeout")
			return false
		}

		var reqHash types.ReqHash
		reqHash.Hash = hash
		res, err := mainClient.QueryTransaction(context.Background(), &reqHash)
		if err != nil {
			fmt.Println("query tx :",err)
			time.Sleep(time.Second)
		}
		if res != nil {
			fmt.Println("res:",res)
			if res.Receipt.Ty == types.ExecOk {
				fmt.Println("ExecOk")
			} else {
				fmt.Println("Exec fail :",res.Receipt.Logs)
			}
			return true
		}
	}
}

func genaddress() (string, crypto.PrivKey) {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		panic(err)
	}
	privto, err := cr.GenKey()
	if err != nil {
		panic(err)
	}
	addrto := address.PubKeyToAddress(privto.PubKey().Bytes())
	return addrto.String(), privto
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName("", types.SECP256K1))
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

func TestCalc(t *testing.T) {
	am := int64(1000000000)
	r := int64(1211)
	day := int64(359)
	amount := big.NewInt(am)
	rate := big.NewInt(r)
	day360 := big.NewInt(int64(360))
	unit := big.NewInt(int64(10000))
	day1 := big.NewInt(day)
	s := big.NewInt(1)
	sm := big.NewInt(1)
	repayson := s.Mul(amount,rate).Mul(s,day1).String()
	repayMath := sm.Mul(sm,day360).Mul(sm,unit).String()
	fmt.Println("son:",repayson," mather:",repayMath)
	repaysonf ,_ := strconv.ParseFloat(repayson,64)
	repayMathf,_ := strconv.ParseFloat(repayMath,64)
	repayfloat64 := repaysonf/repayMathf

	repaycalc := fmt.Sprintf("%0.f",repayfloat64)

	repayint64, _ := strconv.ParseInt(repaycalc, 10, 64)
	fmt.Println("repayfloat64 ",repayfloat64, " repaycalc ",repaycalc," int64 ",repayint64)
}

func round(x float64) int {
	return int(math.Floor(x + 0.5))
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func TestGetAddr(t *testing.T) {
	priv := getprivkey("c03924e244089c22d7322637981a10c9c8b7370e3abb5a1eb21762026a156f1c")
	fmt.Println("addr:",address.PubKeyToAddress(priv.PubKey().Bytes()))
	fmt.Println("addr:",address.ExecAddress("token"))
}

// 检查某一次coins转账交易是否成功
func TestCheckTxSuccess(t *testing.T)  {
	var rpcAddr , txHash string
	rpcAddr = "172.16.100.202:8802"
	txHash = "0x3e592bbb543ac62feb8297fb88a9a4f53d4a9f24848ea4f01244f4903d8173f8"
	 var grpcRecSize int = 30 * 1024 * 1024
	msgRecvOp := grpc.WithMaxMsgSize(grpcRecSize)
	conn, err := grpc.Dial(rpcAddr, grpc.WithInsecure(), msgRecvOp)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := types.NewChain33Client(conn)
	txbyte,_ := common.FromHex(txHash)
	fmt.Println("txbyte:",txbyte)
	res, err := c.QueryTransaction(context.Background(), &types.ReqHash{Hash: txbyte})

	fmt.Println("err",res,err)
	res1, err1 := c.IsSync(context.Background(), &types.ReqNil{})
	fmt.Println("err",res1,err1)

	res2,err2 := c.SendTransaction(context.Background(),&types.Transaction{})
	fmt.Println("err",res2,err2)

}


func TestCheckTx(t *testing.T) {
	txhash := "0x5e26aa27df32218bfc4c709b8d9037acc53f6a6028f783b7af2707b690d94492"
	txbyte,_ := common.FromHex(txhash)

	ok := waitTx(txbyte)
	fmt.Println(ok)
}

func TestOverflow(t *testing.T) {
	var base = int64(100000)
	var rate = int64(28)
	var day = int64(31)
	amount,err := getCalcAmount(base,rate,day,1)
	if err != nil {
		panic(err)
	}
	fmt.Println("repay:",amount)
	var list = [3]int64{30000,30000,40000}
	var total int64
	for k,v := range list {
		repay := (v*amount/base)
		fmt.Println("index:",k," repay :",repay," base :",v)
		total += repay
	}
	fmt.Println("total:",total)
}

func getCalcAmount1(amount,rate,day int64,ty int64) (int64,error) {
	amountbig := big.NewInt(amount)
	ratebig := big.NewInt(rate)
	daybig := big.NewInt(day)
	s := big.NewInt(1)
	if ty == 0 {//年利率

		day360 := big.NewInt(int64(360))
		unit := big.NewInt(int64(10000))
		sm := big.NewInt(1)
		repayson := s.Mul(amountbig,ratebig).Mul(s,daybig).String()
		repaymoth := sm.Mul(sm,day360).Mul(sm,unit).String()
		repaysonf ,err := strconv.ParseFloat(repayson,64)
		if err != nil {
			return 0,err
		}
		repayMathf,err := strconv.ParseFloat(repaymoth,64)
		if err != nil {
			return 0,err
		}
		repayfloat64 := repaysonf/repayMathf

		repaycalc := fmt.Sprintf("%0.f",repayfloat64)

		repayint64, err := strconv.ParseInt(repaycalc, 10, 64)
		if err != nil {
			return 0,err
		}
		return repayint64+amount,nil
	} else {//日利率
		unit := big.NewInt(int64(100000))
		repayson := s.Mul(amountbig,ratebig).Mul(s,daybig).String()
		repaymoth := unit.String()
		repaysonf ,err := strconv.ParseFloat(repayson,64)
		if err != nil {
			return 0,err
		}
		repayMathf,err := strconv.ParseFloat(repaymoth,64)
		if err != nil {
			return 0,err
		}
		repayfloat64 := repaysonf/repayMathf
		repaycalc := fmt.Sprintf("%0.f",repayfloat64)

		repayint64, err := strconv.ParseInt(repaycalc, 10, 64)
		if err != nil {
			return 0,err
		}
		return repayint64+amount,nil
	}

}


func TestCalc1(t *testing.T) {
	day := getSubDays(time.Now().Unix(),int64(1554912000))
	fmt.Println("day:",day)
	amount ,err:= getCalcAmount(int64(100000000),int64(28),day,int64(1))
	if err != nil {
		panic(err)
	}
	fmt.Println("amountL:",amount)

	bigamount := big.NewInt(int64(8000000000))
	bigrepay := big.NewInt(int64(10400000000))
	base := big.NewInt(int64(1))

	repay :=base.Mul(bigamount,bigrepay).Div(base,bigamount).Int64()
	fmt.Println("a:",repay)

	bigtamount := big.NewInt(int64(4000000000))

	temp :=base.Sub(bigamount,bigtamount).Mul(base,bigrepay).Div(base,bigamount).Int64()
	fmt.Println("a:",temp)
}

func TestSignGroupTx(t *testing.T) {
	v1 := &tokennotety.TokennoteCreate{
		Issuer:       issuer,
		IssuerName:issuerName,
		IssuerId:issuerId,
		IssuerPhone:issuerPhone,
		Currency:"SZH00011",
		Acceptor:       issuer,
		AcceptanceDate:1554912000,
		Introduction: tokenIntro,
		Balance:        10000000,
		Rate:        28,
		Repayamount:        0,
	}
	create1 := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionCreate,
		Value: &tokennotety.TokennoteAction_TokennoteCreate{TokennoteCreate: v1},
	}
	tx1 := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(create1),
		Fee:     0,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}

	v2 := &tokennotety.TokennoteCreate{
		Issuer:       to1,
		IssuerName:issuerName,
		IssuerId:issuerId,
		IssuerPhone:issuerPhone,
		Currency:"SZH00012",
		Acceptor:       issuer,
		AcceptanceDate:1554912000,
		Introduction: tokenIntro,
		Balance:        10000000,
		Rate:        28,
		Repayamount:        0,
	}
	create2 := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionCreate,
		Value: &tokennotety.TokennoteAction_TokennoteCreate{TokennoteCreate: v2},
	}
	tx2 := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(create2),
		Fee:     0,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	v3 := &tokennotety.TokennoteCreate{
		Issuer:       to1,
		IssuerName:issuerName,
		IssuerId:issuerId,
		IssuerPhone:issuerPhone,
		Currency:"SZH00013",
		Acceptor:       issuer,
		AcceptanceDate:1554912000,
		Introduction: tokenIntro,
		Balance:        10000000,
		Rate:        28,
		Repayamount:        0,
	}
	create3 := &tokennotety.TokennoteAction{
		Ty:    tokennotety.TokennoteActionCreate,
		Value: &tokennotety.TokennoteAction_TokennoteCreate{TokennoteCreate: v3},
	}
	tx3 := &types.Transaction{
		Execer:  []byte(execName),
		Payload: types.Encode(create3),
		Fee:     0,
		Nonce:   r.Int63(),
		To:      address.ExecAddress(execName),
	}
	txs := []*types.Transaction{}
	txs = append(txs,tx1,tx2,tx3)
	txgroup,_ := types.CreateTxGroup(txs)
	txgroup.SignN(0,types.SECP256K1, privkey)
	txgroup.SignN(1,types.SECP256K1, to1Privkey)
	txgroup.SignN(1,types.SECP256K1, to1Privkey)

	tx := txgroup.Tx()
	tx.Sign(types.SECP256K1, privkey)
	fmt.Println("grouptx:",hex.EncodeToString(types.Encode(tx)))
	fmt.Println("size:",types.Size(tx1))
	fmt.Println("size:",types.Size(tx2))
}
func TestSignGroupTx1(t *testing.T) {
	//tx := &types.CreateTx{
	//	To:szh,
	//	Amount:1e6,
	//	Note:[]byte("test"),
	//	Fee:1e5,
	//
	//}
	//hex,err := paraClient.CreateRawTransaction(context.Background(),tx)
}

func TestGetBalance(t *testing.T) {

	res, err := mainClient.GetBalance(context.Background(),&types.ReqBalance{
		Addresses:[]string{issuer},
		Execer:"coins",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("res",res)

	fmt.Println("address:",address.ExecAddress("ticket"))
}


func TestNoneTx(t *testing.T) {
	tx := &types.Transaction{
		Execer:  []byte("user.p.sy.none"),
		Payload: []byte("test"),
		Fee:     feeForToken,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("user.p.sy.none"),
	}
	tx.Sign(types.SECP256K1, privkey)

	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
}

func TestLocalKey(t *testing.T) {
	minkey := len([]byte("LODB"))+len([]byte("tokennote")) + 2
	key := []byte("LODB-tokennote-createdNotes")
	fmt.Println([]byte("-"))
	if key[14] == '-' {
		fmt.Println("ok",key[14],minkey)
	} else {
		fmt.Println("fail")
	}
	pub,_ := hex.DecodeString("0236813da2fced52ca344ab5e062b4f88fa5783c7f688dfc7ab4d625c4b6c06a22")
	fmt.Println("addr",address.PubKeyToAddress(pub).String(),address.PubKeyToAddr(pub))

}


func TestGBTSign(t *testing.T) {
	txhex := fmt.Sprintf("%018d", 999*100000 + 0)
	fmt.Println(txhex)
	s := "0a15757365722e702e676274636861696e2e636f696e73122e18010a2a1080e788d8032222314a354556464b6571445755776e6e584c556576534269396347363761584c677a791a6e08011221034a5408ba67b7b3ffafb3a840c168fddd15eda9a4caa9c1dc50612f45c502145c1a473045022100e9fa32ec5fd1f3df955a6e3f0dcb6951b5957cc8ac80dc9449cd5d050c24e76d02207c94521017e07911643bfa5aee2ac220aeb8e263fdb269fb7b19b0c5bf0c658320a08d062880c6868f0130f7ccc39ee991f6d31c3a22314737786570574771746231376f4e7772364d7a796f626a766d4534636b4853464e"
	bs ,_ := common.FromHex(s)
	var tx types.Transaction
	err := types.Decode(bs,&tx)
	if err != nil {
		panic(err)
	}
	fmt.Println(tx,address.PubKeyToAddr(tx.Signature.Pubkey),tx.Fee,string(tx.Execer))
	var payload cty.CoinsAction
	_ = types.Decode(tx.Payload,&payload)
	switch payload.Ty {
	case cty.CoinsActionTransfer:
		fmt.Println("transfer",payload.GetTransfer().Amount,payload.GetTransfer().To,payload.GetTransfer().Cointoken)
	case cty.CoinsActionTransferToExec:
		fmt.Println("toExec")
	case cty.CoinsActionWithdraw:
		fmt.Println("withdraw")
	default:
		panic("wrong")
	}

}