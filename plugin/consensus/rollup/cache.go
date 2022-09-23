package rollup

import (
	"bytes"
	"sync"

	rtypes "github.com/33cn/plugin/plugin/dapp/rollup/types"
)

const (
	// 缓存已提交的历史checkpoint数量(主要应对主链可能的分叉回滚情况)
	historyCacheCount = 10
)

type commitCache struct {
	lock          sync.RWMutex
	minCacheRound int64
	cpList        map[int64]*rtypes.CheckPoint
	signList      map[int64]*validatorSignMsgSet
}

func newCommitCache(currRound int64) *commitCache {

	c := &commitCache{minCacheRound: currRound}
	c.cpList = make(map[int64]*rtypes.CheckPoint, 32)
	c.signList = make(map[int64]*validatorSignMsgSet, 32)
	return c
}

func (c *commitCache) addCheckPoint(batch *rtypes.CheckPoint) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.cpList[batch.CommitRound] = batch
}

func (c *commitCache) getCheckPoint(round int64) *rtypes.CheckPoint {

	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.cpList[round]
}

func (c *commitCache) addValidatorSign(isSelf bool, sign *rtypes.ValidatorSignMsg) {
	c.lock.Lock()
	defer c.lock.Unlock()
	set, ok := c.signList[sign.CommitRound]
	if !ok {
		set = &validatorSignMsgSet{others: make([]*rtypes.ValidatorSignMsg, 0, 8)}
		c.signList[sign.CommitRound] = set
	}
	if isSelf {
		set.self = sign
		return
	}

	// 检测是否有重复
	for _, other := range set.others {
		if bytes.Equal(sign.PubKey, other.PubKey) {
			return
		}
	}
	set.others = append(set.others, sign)
}

// clean already commit batch and sign
func (c *commitCache) cleanHistory(currRound int64) {
	c.lock.Lock()
	defer c.lock.Unlock()

	maxDelRound := currRound - historyCacheCount
	for i := c.minCacheRound; i <= maxDelRound; i++ {
		delete(c.signList, i)
		delete(c.cpList, i)
	}

	if maxDelRound >= c.minCacheRound {
		c.minCacheRound = maxDelRound + 1
	}
}

//

func (c *commitCache) getPreparedCheckPoint(round int64, aggreSign aggreSignFunc) *rtypes.CheckPoint {

	c.lock.RLock()
	defer c.lock.RUnlock()

	signSet := c.signList[round]
	pubs, aSign := aggreSign(signSet)
	if pubs == nil {
		return nil
	}
	cp := c.cpList[round]

	cp.ValidatorPubs = pubs
	cp.AggregateValidatorSign = aSign
	return cp
}
