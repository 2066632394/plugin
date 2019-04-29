package executor

import (
	"testing"
	"github.com/33cn/chain33/types"
	gtypes "github.com/33cn/plugin/plugin/dapp/goldbill/types"
	"math/rand"
	"google.golang.org/grpc"
	"time"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/crypto"
	"fmt"
	"encoding/hex"
	"context"
)

var (
	mainNetgrpcAddr = "172.16.100.55:8802"
	mainClient types.Chain33Client
	r *rand.Rand
	priv = getprivkey("0x016959665e2dd245ffef62a73df1dfe633308e57421366e260da136731137961")
)

func init() {

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	maxReceLimit := grpc.WithMaxMsgSize(30*1024*1024)
	conn, err := grpc.Dial(mainNetgrpcAddr, grpc.WithInsecure(),maxReceLimit)
	if err != nil {
		panic(err)
	}
	mainClient = types.NewChain33Client(conn)
}

func TestRegisterUser (t *testing.T) {
	v := &gtypes.GoldbillAction_RegisterUser{
		RegisterUser: &gtypes.GoldbillRegisterUser{
			UserId: "SZH000001",
			UserPubkey:priv.PubKey().Bytes(),
			UserType: gtypes.GoldbillUserType_UT_USER,
		}}
	mint := &gtypes.GoldbillAction{Value: v, Ty: gtypes.GoldbillActionType_RegisterUser}

	tx := &types.Transaction{
		Execer:  []byte("goldbill"),
		Payload: types.Encode(mint),
		Fee:     1e5,
		Nonce:   r.Int63(),
		To:      address.ExecAddress("goldbill"),
	}
	tx.Nonce = r.Int63()
	tx.Sign(types.SECP256K1, priv)
	reply, err := mainClient.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)

		t.Error(err)
		return
	}
	fmt.Println("replay:",reply,common.ToHex(reply.Msg))

	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error(err)
		return
	}
	fmt.Println("txhash:",hex.EncodeToString(tx.Hash()))
	if !waitTx(tx.Hash()) {
		t.Error(err)
		return
	}
	time.Sleep(5 * time.Second)
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

// ExecTypeGet  获取类型值
type execTypeGet interface {
	GetTy() int32
}