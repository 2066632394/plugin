package executor

import (
	"fmt"
	dbm "github.com/33cn/chain33/common/db"
)
var (
	orderPrefix = "mavl-" + "presale" + "-"
	adminPubkey = "ba2e9b95cfd8348b1e147f5cbeb55a921a7fbd10e091825bb6a4cf826aad06e4"

)

func proKeyMarketId(marketId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s+%d", "market", marketId))
}

func proKeyMarketIdJudge(id int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"judgeMarket"+"%d", id))
}

func proKeyUserId(uid int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s+%d", "user", uid))
}

func proKeyUidJudge(uid int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"judgeUid"+"%d", uid))
}

func proKeyPdId(pdId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"pd", "%d", pdId))
}

func proKeyOrderId(orderId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"order", "%d", orderId))
}

func proKeyToken(token string, marketId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s_%d", token, marketId))
}

func calcOrderKey(hash string) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s", hash))
}

func proKeyAddress(uid int64, addressId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%d+%d", uid, addressId))
}

func proKeyBTY(bty string, marketId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s_%d", bty, marketId))
}

func proKeyAdminReg(adminId int64) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%d", adminId))
}

func KeyMarketReg(hash []byte) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s", hash))
}

func proString(chuan string) []byte {
	return []byte(fmt.Sprintf(orderPrefix+"%s", chuan))
}

type order struct {
}

func (o *order) save(db dbm.KV, key, value []byte) {
	db.Set([]byte(key), value)
}