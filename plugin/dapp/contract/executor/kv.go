package executor

import (
	"fmt"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

var (
	contractPrefix = fmt.Sprintf("mavl-%s-contract-", cty.ContractX)
	userPrefix     = fmt.Sprintf("mavl-%s-user-", "gyl")
)

func calcContractKey(contractId string) []byte {
	return []byte(fmt.Sprintf(contractPrefix+"%s", contractId))
}

func calcUserKey(userId string) []byte {
	return []byte(fmt.Sprintf(userPrefix+"%s", userId))
}
