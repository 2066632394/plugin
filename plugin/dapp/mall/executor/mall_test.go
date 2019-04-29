package executor

import (
	"fmt"
	"google.golang.org/grpc"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/common"
	"golang.org/x/net/context"
	"testing"
)

func TestGrpc(t *testing.T){

	con,err := grpc.Dial("172.16.100.148:8802",grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := types.NewChain33Client(con)
	h := "0x6a7375cb6fe885aba172c917fd9ee31793c0f90bae0a5fe41acfd58d09cbc80e"
	hash,err := common.FromHex(h)
	if err != nil {
		panic(err)
	}

	res,err := c.QueryTransaction(context.Background(),&types.ReqHash{Hash:hash})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestGetLogType(t *testing.T){
	res := types.GetLogName([]byte("mall"),1)
	fmt.Println("logtype:",res)
	logMap := types.LoadExecutorType("mall")
	fmt.Println(logMap.GetLogMap())

}