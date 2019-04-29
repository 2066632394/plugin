package types

import (
	"testing"
	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common/address"
	"fmt"
	//"github.com/33cn/chain33/rpc/jsonclient"
	"encoding/json"
	"encoding/hex"
	"time"
	"math"
)

func TestSubTime(t *testing.T) {
	s := int64(1546481093)

	t1 := time.Unix(s,0)
	time.Sleep(time.Second*10)
	s1 := int64(1548900294)
	t2 := time.Unix(s1,0)
	fmt.Println(int64(math.Ceil(t2.Sub(t1).Hours()/24)))

}

func TestNumber(t *testing.T) {
	str := "2XASD@"

	sb := []byte(str)
	for k,a := range sb {
		res := (a <= 'Z' && a >= 'A')
		res1 := (a <= '9' && a >= '0')
		fmt.Println(k,":",a,":",res,":",res1)
	}


}

func TestGetAddr(t *testing.T) {
	pub := "02eee16a95638b765516bb63e29649820221783d94fe68344b798e5547fe8eeda0"
	pubbyte,err := common.FromHex(pub)
	if err != nil {
		panic(err)
	}
	addr := address.PubKeyToAddress(pubbyte).String()
	fmt.Println("addr:",addr)
}
func TestSign(t *testing.T) {
	privstr := "02763d0a06742ff9e68398f2961d5586c0c05ca91b232ff0334e4eda1354b13af7"
	key,err := common.FromHex(privstr)
	if err != nil {
		panic(err)
	}
	cr, err := crypto.New(types.GetSignName("", 1))
	if err != nil {
		panic(err)
	}
	k, err := cr.PrivKeyFromBytes(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("k:",k)


}

func TestMallUserReg(t *testing.T) {
	privkey := getprivkey("0x43e3c5624c378b3c9027cc1e8c4d52521432e59c4a71bc963dc63f9846b08dcb")
	userreg := &MallUserRegister{Uid:"test",Phone:13754338419}
	tx := &types.Transaction{
		Execer:[]byte("user.p.malltest.mall"),
		To:address.ExecAddress("user.p.malltest.mall"),
		Fee:100000,
		Payload:types.Encode(&MallAction{Value:&MallAction_MallUserRegister{userreg},Ty:1}),
	}
	txx := types.Encode(tx)
	data := hex.EncodeToString(txx)
	fmt.Println("hexstring:",data,"pub:",hex.EncodeToString(privkey.PubKey().Bytes()))
	//mocker := testnode.New("--free--", nil)
	//client,err := jsonclient.NewJSONClient("http://172.16.100.49:8801")
	//if err != nil {
	//	panic(err)
	//}
	//var result interface{}
	//var notx types.NoBalanceTx
	//notx.Expire = "1h"
	//notx.PayAddr = "15M1iwuk5eGbqngsjUbzNosTjdiVq7ryHc"
	//notx.TxHex = data
	//err = client.Call("CreateNoBalanceTransaction",&notx,result)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("result:",result)

}

type clientResponse struct {
	Id     uint64           `json:"id"`
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
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
