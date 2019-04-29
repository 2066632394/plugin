package executor

import (
	"context"
	"fmt"
	"github.com/33cn/chain33/common"
	"google.golang.org/grpc"
	"math/rand"
	"testing"
	"time"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

const (
	FEE = 0
)

var (
	ContractAddr    string
	conn            *grpc.ClientConn
	r               *rand.Rand
	c               types.Chain33Client
	cfg             *types.Config
	superPriv       crypto.PrivKey
	company1Privkey crypto.PrivKey
	company2Privkey crypto.PrivKey
	company3Privkey crypto.PrivKey
	localIp         = "127.0.0.1:8802"
	onlineIp        = "172.16.100.71:8802"
)

func init() {
	conn, err := grpc.Dial(localIp, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c = types.NewChain33Client(conn)
	r = rand.New(rand.NewSource(types.Now().UnixNano()))
	ContractAddr = address.ExecAddress(cty.ContractX)
	superPriv = getprivkey("9a188e0994b04868504913e6f433202a781d21f02c842100750815548e8849b44be3ca87cab28f156516a65cb91acf1d33020f3c4214005dd95d44ff872062d3")
	company1Privkey = getprivkey("786a6c88adcfbc0cfc51c29149ded4e8d3ad65159f2d403554fa684b0f0764d5cb381440b02e9e7befab2b6fe24de3ca30a39d1ad87399bad0c311dfb58b851c")
	company2Privkey = getprivkey("b0c5f1d6ad55053500c2fa085feef1d5719402e7c2494dc48529c37f736139d05ee9699538e3c6e09436c93d8e066ee2a0b43bfc8e24b58946444d5eab08834c")
	company3Privkey = getprivkey("6f2697d4c1ee5685cdce6fe10b91d12c9da5a82c040a2e69c4b1cb3cba910dad9db8effb23e787867f54a00c688c02edf8e7e14c3a6b1578542a5f91617134d1")
}

func TestAddr(t *testing.T) {
	t.Log(ContractAddr)
}

func TestActions(t *testing.T) {
	fmt.Println("\nTestActions start")
	defer fmt.Println("TestActions end")
	//ProcessArray := []int{1, 2, 3} //注册用户
	//ProcessArray := []int{4, 5} //创建合同后撤销
	ProcessArray := []int{4, 6, 7, 8} //创建合同后签名后拒签
	var tx *types.Transaction
	for _, v := range ProcessArray {
		switch v {
		//case 1: //注册公司1
		//	tx = establishTx(establishRegisterCompany1Act(), superPriv)
		//case 2: //注册公司2
		//	tx = establishTx(establishRegisterCompany2Act(), superPriv)
		//case 3: //注册公司3
		//	tx = establishTx(establishRegisterCompany3Act(), superPriv)
		case 4: //创建合同
			tx = establishTx(establishCreateAct(), company1Privkey)
		case 5: //撤销合同
			tx = establishTx(establishCancelAct(), company1Privkey)
		case 6: //修改合同
			tx = establishTx(establishModifyAct(), company1Privkey)
		case 7: //签署合同
			tx = establishTx(establishSignAct(), company2Privkey)
		case 8: //拒绝合同
			tx = establishTx(establishRejectAct(), company3Privkey)
		default:

		}
		fmt.Println(common.Bytes2Hex(tx.GetSignature().GetPubkey()))

		reply, err := c.SendTransaction(context.Background(), tx)
		if err != nil {
			fmt.Println("err", err)
			t.Error(err)
			return
		}
		if !reply.IsOk {
			fmt.Println("err = ", reply.GetMsg())
			t.Error("test error")
			return
		}
		fmt.Printf("ty:%d hash:%s\n", v, common.Bytes2Hex(reply.GetMsg()))
		//if !waitTx(tx.Hash()){
		//	t.Error(ErrTest)
		//	return
		//}
		time.Sleep(100 * time.Millisecond)
	}
}

func TestOne(t *testing.T) {
	fmt.Println("\nTestOne start")
	defer fmt.Println("TestOne end")
	var tx *types.Transaction
	v := 8
	switch v {
	//case 1: //注册公司1
	//	tx = establishTx(establishRegisterCompany1Act(), superPriv)
	//case 2: //注册公司2
	//	tx = establishTx(establishRegisterCompany2Act(), superPriv)
	//case 3: //注册公司3
	//	tx = establishTx(establishRegisterCompany3Act(), superPriv)
	case 4: //创建合同
		tx = establishTx(establishCreateAct(), company1Privkey)
	case 5: //撤销合同
		tx = establishTx(establishCancelAct(), company1Privkey)
	case 6: //修改合同
		tx = establishTx(establishModifyAct(), company1Privkey)
	case 7: //签署合同
		tx = establishTx(establishSignAct(), company2Privkey)
	case 8: //拒绝合同
		tx = establishTx(establishRejectAct(), company3Privkey)
	default:

	}
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Println("err", err)
		t.Error(err)
		return
	}
	if !reply.IsOk {
		fmt.Println("err = ", reply.GetMsg())
		t.Error("test error")
		return
	}
	fmt.Printf("ty:%d hash:%s\n", v, common.Bytes2Hex(reply.GetMsg()))
}

func TestDuplicate(t *testing.T) {
	t.Log(isDuplicate("1", "1", "3"))
}

func establishTx(contractAction *cty.ContractAction, privKey crypto.PrivKey) *types.Transaction {

	tx := &types.Transaction{
		Execer:  cty.ExecerContract,
		Payload: types.Encode(contractAction),
		Fee:     FEE,
		To:      ContractAddr,
		Nonce:   r.Int63(),
	}
	tx.Sign(types.ED25519, privKey)

	return tx
}

func establishCreateAct() *cty.ContractAction {
	contractCreate := &cty.ContractAction_Create{
		Create: &cty.ContractCreate{
			OperatorId:         "1",
			ContractId:         "11",
			ContractHash:       "hdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdius",
			ContractName:       "1创建的1",
			Amount:             100,
			SignatoryIds:       []string{"2", "3"},
			IsDraft:            true,
			SignedContractHash: "hdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdaadiushdahdiushdahdiushdahdius",
			OperateTime:        uint64(time.Now().Unix()),
		},
	}
	return &cty.ContractAction{
		Value: contractCreate,
		Ty:    cty.ContractActionCreate,
	}
}

func establishCancelAct() *cty.ContractAction {
	contractCancel := &cty.ContractAction_Cancel{
		Cancel: &cty.ContractCancel{
			OperatorId:  "1",
			ContractId:  "11",
			OperateTime: uint64(time.Now().Unix()),
		},
	}
	return &cty.ContractAction{
		Value: contractCancel,
		Ty:    cty.ContractActionCancel,
	}
}

func establishModifyAct() *cty.ContractAction {
	contractModify := &cty.ContractAction_Modify{
		Modify: &cty.ContractModify{
			OperatorId: "1",
			ContractId: "11",
			//ContractHash: "hdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdius",
			//ContractName: "1创建的1",
			Amount: 200,
			//SignatoryIds: []string{"2", "3"},
			SignedContractHash: "hqwhdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdaadiushdahdiushdahdiushdahdius",
			OperateTime:        uint64(time.Now().Unix()),
		},
	}
	return &cty.ContractAction{
		Value: contractModify,
		Ty:    cty.ContractActionModify,
	}
}

func establishSignAct() *cty.ContractAction {
	contractSign := &cty.ContractAction_Sign{
		Sign: &cty.ContractSign{
			OperatorId:         "2",
			ContractId:         "11",
			SignedContractHash: "hqwhddsahdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdahdiushdaadiushdahdiushdahdiushdahdius",
			OperateTime:        uint64(time.Now().Unix()),
		},
	}
	return &cty.ContractAction{
		Value: contractSign,
		Ty:    cty.ContractActionSign,
	}
}

func establishRejectAct() *cty.ContractAction {
	contractReject := &cty.ContractAction_Reject{
		Reject: &cty.ContractReject{
			OperatorId:  "3",
			ContractId:  "11",
			OperateTime: uint64(time.Now().Unix()),
		},
	}
	return &cty.ContractAction{
		Value: contractReject,
		Ty:    cty.ContractActionReject,
	}
}

//func establishRegisterCompany1Act() *cty.ContractAction {
//	contractRegister := &cty.ContractAction_Register{
//		Register: &cty.ContractRegister{
//			Id:     "1",
//			Name:   "公司1",
//			PubKey: company1Privkey.PubKey().KeyString(),
//		},
//	}
//	return &cty.ContractAction{
//		Value: contractRegister,
//		Ty:    cty.ContractActionRegister,
//	}
//}
//
//func establishRegisterCompany2Act() *cty.ContractAction {
//	contractRegister := &cty.ContractAction_Register{
//		Register: &cty.ContractRegister{
//			Id:     "2",
//			Name:   "公司2",
//			PubKey: company2Privkey.PubKey().KeyString(),
//		},
//	}
//	return &cty.ContractAction{
//		Value: contractRegister,
//		Ty:    cty.ContractActionRegister,
//	}
//}
//
//func establishRegisterCompany3Act() *cty.ContractAction {
//	contractRegister := &cty.ContractAction_Register{
//		Register: &cty.ContractRegister{
//			Id:     "3",
//			Name:   "公司3",
//			PubKey: company3Privkey.PubKey().KeyString(),
//		},
//	}
//	return &cty.ContractAction{
//		Value: contractRegister,
//		Ty:    cty.ContractActionRegister,
//	}
//}

//生成私钥
func genAddress() (string, crypto.PrivKey) {
	cr, err := crypto.New(types.GetSignName(cty.ContractX, types.ED25519))
	if err != nil {
		panic(err)
	}
	priTo, err := cr.GenKey()
	if err != nil {
		panic(err)
	}
	addrTo := address.PubKeyToAddress(priTo.PubKey().Bytes())
	return addrTo.String(), priTo
}

func getprivkey(key string) crypto.PrivKey {
	cr, err := crypto.New(types.GetSignName(cty.ContractX, types.ED25519))
	if err != nil {
		panic(err)
	}
	bkey, err := common.FromHex(key)
	//fmt.Println(bkey)
	if err != nil {
		panic(err)
	}
	priv, err := cr.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}
