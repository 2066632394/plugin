package executor

import (
	"fmt"
)


var (
	fzmsupplyPrefix = "mavl-fzmsupply-"
	fzmsupplyHistory = "fzmsupplynfccodehistory-"
	fzmsupplyTxsCount = "fzmsupplytxscount-"
)


func fzmsupplyKeyUser(id string) []byte {
	return []byte( fzmsupplyPrefix+  "user-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyAdmin(id string) []byte {
	return []byte( fzmsupplyPrefix+  "admin-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyOfficial(id string) []byte {
	return []byte( fzmsupplyPrefix+  "official-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyRole(id string) []byte {
	return []byte( fzmsupplyPrefix+  "role-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyContract(id string) []byte {
	return []byte( fzmsupplyPrefix+  "contract-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyPayment(id string) []byte {
	return []byte( fzmsupplyPrefix+  "payment-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyLoan(id string) []byte {
	return []byte( fzmsupplyPrefix+  "loan-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyAsset(id string) []byte {
	return []byte( fzmsupplyPrefix+  "asset-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyDeposit(id string) []byte {
	return []byte( fzmsupplyPrefix+  "deposit-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyWithdraw(id string) []byte {
	return []byte( fzmsupplyPrefix+  "withdraw-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyBatch(id string) []byte {
	return []byte( fzmsupplyPrefix+  "batch-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyFinance(id string) []byte {
	return []byte( fzmsupplyPrefix+  "finance-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyMortgage(id string) []byte {
	return []byte( fzmsupplyPrefix+  "mortgage-" + fmt.Sprintf("%s", id))
}

func fzmsupplyKeyPlatform() []byte {
	return []byte( fzmsupplyPrefix+  "platform")
}