package executor

import (
	"testing"
	h "github.com/33cn/plugin/plugin/dapp/hacker/types"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"math/rand"
	"fmt"
	"encoding/hex"
	"time"
)

func TestAddBill(t *testing.T) {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	privkey := getprivkey("0x016959665e2dd245ffef62a73df1dfe633308e57421366e260da136731137961")
	v := &h.HackerAddBill{
		StockNumber: "cd20190226003",
		StockName:"花生油",
		Brand:"金龙鱼",
		BatchRequest:"IXXXX9001-2019",
		PledgeRate:10.1,
		BasicUnit: "箱",
		CommodityCode:"1580612",
		ExpirationDate: int64(1551181881),
		PledgePrice:   100000.0,
		EarlyWarningDate:int64(1551181881),
		Specification:"ddd",
	}
	create := &h.HackerAction{
		Ty:    h.HackerAddBillAction,
		Value: &h.HackerAction_AddBill{ v},
	}
	tx := &types.Transaction{
		Execer:  []byte("hacker"),
		Payload: types.Encode(create),
		Fee:     0,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("hacker"),
	}
	tx.Sign(types.SECP256K1, privkey)
	fmt.Println("tx:",hex.EncodeToString(types.Encode(tx)))
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