package executor

import (
	"github.com/33cn/chain33/account"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	cty "github.com/33cn/plugin/plugin/dapp/contract/types"
)

type DB struct {
	contract cty.Contract
	//MODIFY: user     cty.User
}

func newDB() *DB {
	return &DB{}
}

func (c *DB) save(db dbm.KV) []*types.KeyValue {
	set := c.getKVSet()
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}

	return set
}

func (c *DB) getKVSet() (kvSet []*types.KeyValue) {
	if c.contract.ContractId != "" {
		kvSet = append(kvSet, &types.KeyValue{Key: calcContractKey(c.contract.ContractId), Value: types.Encode(&c.contract)})
	}
	//MODIFY: if c.user.Id != "" {
	//MODIFY: 	kvSet = append(kvSet, &types.KeyValue{Key: calcUserKey(c.user.Id), Value: types.Encode(&c.user)})
	//MODIFY: }

	return kvSet
}

type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txHash       []byte
	fromAddr     string
	blockTime    int64
	height       int64
	execAddr     string
}

func newContractAction(r *Contract, tx *types.Transaction) *Action {
	hash := tx.Hash()
	fromAddr := tx.From()
	return &Action{r.GetCoinsAccount(), r.GetStateDB(), hash,
		fromAddr, r.GetBlockTime(), r.GetHeight(), dapp.ExecAddress(r.GetName())}
}

func (action *Action) contractCreate(create *cty.ContractCreate) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	c := newDB()

	//数据格式验证
	if err := create.Validate(); err != nil {
		return nil, cty.ErrFormat
	}

	//判断合同是否已存在
	_, err := action.getContractById(create.GetContractId())
	if err != types.ErrNotFound {
		if err != nil {
			return nil, err
		}
		return nil, cty.ErrContractExists
	}
	//判断用户id是否重复
	if isDuplicate(create.GetContractId(), create.GetSignatoryIds()) {
		return nil, cty.ErrDuplicateUserId
	}
	//新建合同
	contract := &cty.Contract{
		ContractId:   create.GetContractId(),
		ContractHash: create.GetContractHash(),
		ContractName: create.GetContractName(),
		Amount:       create.GetAmount(),
		CreateTime:   create.GetOperateTime(),
		UpdateTime:   create.GetOperateTime(),
	}
	_, err = action.getUserById(create.GetOperatorId())
	if err != nil {
		return nil, err
	}
	//根据是否为草稿设置不同的状态
	if create.GetIsDraft() {
		contract.Creator = &cty.Creator{
			Id: create.GetOperatorId(),
			//MODIFY: Name:       creator.GetName(),
			SignStatus: cty.SignStatus_SS_Unsigned,
		}
		contract.Status = cty.ContractStatus_CS_Draft
	} else {
		if len(create.GetSignedContractHash()) != 128 {
			return nil, cty.ErrHash
		}
		if create.GetContractHash() == create.GetSignedContractHash() {
			return nil, cty.ErrSameHash
		}
		contract.Creator = &cty.Creator{
			Id:           create.GetOperatorId(),
			ContractHash: create.GetContractHash(),
			//MODIFY: Name:         creator.GetName(),
			SignTime:   create.GetOperateTime(),
			SignStatus: cty.SignStatus_SS_Signed,
		}
		contract.Status = cty.ContractStatus_CS_Ineffective
	}
	//查找所有的签署者是否已存在并写入合同
	var signatories []*cty.Signatory
	for _, v := range create.GetSignatoryIds() {
		_, err := action.getUserById(v)
		if err != nil {
			return nil, err
		}
		signatories = append(signatories, &cty.Signatory{
			Id: v,
			//MODIFY: Name:       signatory.GetName(),
			SignStatus: cty.SignStatus_SS_Unsigned,
		})
	}
	contract.Signatories = signatories

	c.contract = *contract
	kv = append(kv, c.save(action.db)...)
	logs = append(logs, &types.ReceiptLog{Ty: cty.TyLogCreateContract, Log: types.Encode(contract)})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (action *Action) contractCancel(cancel *cty.ContractCancel) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	c := newDB()

	//数据格式验证
	if err := cancel.Validate(); err != nil {
		return nil, cty.ErrFormat
	}

	//判断合同是否已存在，并验证签名状态等信息
	contract, err := action.getContractById(cancel.GetContractId())
	if err != nil {
		return nil, err
	}
	if contract.GetCreator().GetId() != cancel.GetOperatorId() {
		return nil, cty.ErrPermissionDenied
	}
	if contract.GetStatus() != cty.ContractStatus_CS_Ineffective && contract.GetStatus() != cty.ContractStatus_CS_Draft {
		return nil, cty.ErrContractStatus
	}
	if cancel.GetOperateTime() < contract.GetUpdateTime() {
		return nil, cty.ErrOperateTime
	}
	for _, v := range contract.GetSignatories() {
		if v.GetSignStatus() != cty.SignStatus_SS_Unsigned {
			return nil, cty.ErrSignStatus
		}
	}
	contract.Status = cty.ContractStatus_CS_Canceled
	contract.Creator.CancelTime = cancel.GetOperateTime()
	contract.UpdateTime = cancel.GetOperateTime()

	c.contract = *contract
	kv = append(kv, c.save(action.db)...)
	logs = append(logs, &types.ReceiptLog{Ty: cty.TyLogCancelContract, Log: types.Encode(contract)})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (action *Action) contractModify(modify *cty.ContractModify) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	c := newDB()

	//数据格式验证
	if err := modify.Validate(); err != nil {
		return nil, cty.ErrFormat
	}

	//判断合同是否已存在，并验证签名状态等信息
	contract, err := action.getContractById(modify.GetContractId())
	if err != nil {
		return nil, err
	}
	if contract.GetCreator().GetId() != modify.GetOperatorId() {
		return nil, cty.ErrPermissionDenied
	}
	if contract.GetStatus() != cty.ContractStatus_CS_Ineffective && contract.GetStatus() != cty.ContractStatus_CS_Draft {
		return nil, cty.ErrContractStatus
	}
	if modify.GetOperateTime() < contract.GetUpdateTime() {
		return nil, cty.ErrOperateTime
	}
	for _, v := range contract.GetSignatories() {
		if v.GetSignStatus() != cty.SignStatus_SS_Unsigned {
			return nil, cty.ErrSignStatus
		}
	}
	//修改合同
	if modify.GetContractHash() != "" {
		if len(modify.GetContractHash()) != 128 {
			return nil, cty.ErrHash
		}
		contract.ContractHash = modify.GetContractHash()
		if contract.GetStatus() != cty.ContractStatus_CS_Draft && modify.GetSignedContractHash() == "" {
			return nil, cty.ErrSignedHash
		}
	}
	if modify.GetSignedContractHash() != "" && contract.GetCreator().GetSignStatus() == cty.SignStatus_SS_Signed {
		if contract.GetContractHash() == modify.GetSignedContractHash() {
			return nil, cty.ErrSameHash
		}
		contract.Creator.ContractHash = modify.GetSignedContractHash()
		contract.Creator.SignTime = modify.GetOperateTime()
	}
	if modify.GetContractName() != "" {
		contract.ContractName = modify.GetContractName()
	}
	if modify.GetAmount() != 0 {
		contract.Amount = modify.GetAmount()
	}
	if len(modify.GetSignatoryIds()) != 0 {
		//判断用户id是否重复
		if isDuplicate(modify.GetContractId(), modify.GetSignatoryIds()) {
			return nil, cty.ErrDuplicateUserId
		}
		var signatories []*cty.Signatory
		for _, v := range modify.GetSignatoryIds() {
			_, err := action.getUserById(v)
			if err != nil {
				return nil, err
			}
			signatories = append(signatories, &cty.Signatory{
				Id: v,
				//MODIFY: Name:       signatory.GetName(),
				SignStatus: cty.SignStatus_SS_Unsigned,
			})
		}
		contract.Signatories = signatories
	}
	contract.UpdateTime = modify.GetOperateTime()

	c.contract = *contract
	kv = append(kv, c.save(action.db)...)
	logs = append(logs, &types.ReceiptLog{Ty: cty.TyLogModifyContract, Log: types.Encode(contract)})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (action *Action) contractSign(sign *cty.ContractSign) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	c := newDB()

	//数据格式验证
	if err := sign.Validate(); err != nil {
		return nil, cty.ErrFormat
	}

	//判断合同是否已存在，并验证签名状态等信息
	contract, err := action.getContractById(sign.GetContractId())
	if err != nil {
		return nil, err
	}
	if contract.GetStatus() != cty.ContractStatus_CS_Ineffective && contract.GetStatus() != cty.ContractStatus_CS_Draft {
		return nil, cty.ErrContractStatus
	}
	if sign.GetOperateTime() < contract.GetUpdateTime() {
		return nil, cty.ErrOperateTime
	}
	//判断签名者为
	if contract.GetCreator().GetId() == sign.GetOperatorId() {
		if contract.GetCreator().GetSignStatus() != cty.SignStatus_SS_Unsigned {
			return nil, cty.ErrSignStatus
		}
		if contract.GetContractHash() == sign.GetSignedContractHash() {
			return nil, cty.ErrSameHash
		}
		contract.Creator.SignStatus = cty.SignStatus_SS_Signed
		contract.Creator.ContractHash = sign.GetSignedContractHash()
		contract.Creator.SignTime = sign.GetOperateTime()
	} else {
		signatoryIndex := -1
		for i, v := range contract.GetSignatories() {
			if v.GetId() == sign.GetOperatorId() {
				signatoryIndex = i
				break
			}
		}
		if signatoryIndex < 0 {
			return nil, cty.ErrNotInSignatories
		}
		if contract.GetSignatories()[signatoryIndex].GetSignStatus() != cty.SignStatus_SS_Unsigned {
			return nil, cty.ErrSignStatus
		}
		if contract.GetContractHash() == sign.GetSignedContractHash() {
			return nil, cty.ErrSameHash
		}
		contract.Signatories[signatoryIndex].SignStatus = cty.SignStatus_SS_Signed
		contract.Signatories[signatoryIndex].ContractHash = sign.GetSignedContractHash()
		contract.Signatories[signatoryIndex].SignTime = sign.GetOperateTime()
	}
	if isContractEffective(contract) {
		contract.Status = cty.ContractStatus_CS_Effective
	} else {
		contract.Status = cty.ContractStatus_CS_Ineffective
	}
	contract.UpdateTime = sign.GetOperateTime()

	c.contract = *contract
	kv = append(kv, c.save(action.db)...)
	logs = append(logs, &types.ReceiptLog{Ty: cty.TyLogSignContract, Log: types.Encode(contract)})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func (action *Action) contractReject(reject *cty.ContractReject) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	c := newDB()

	//数据格式验证
	if err := reject.Validate(); err != nil {
		return nil, cty.ErrFormat
	}

	//判断合同是否已存在，并验证签名状态等信息
	contract, err := action.getContractById(reject.GetContractId())
	if err != nil {
		return nil, err
	}
	if contract.GetStatus() != cty.ContractStatus_CS_Ineffective && contract.GetStatus() != cty.ContractStatus_CS_Draft {
		return nil, cty.ErrContractStatus
	}
	if reject.GetOperateTime() < contract.GetUpdateTime() {
		return nil, cty.ErrOperateTime
	}
	signatoryIndex := -1
	for i, v := range contract.GetSignatories() {
		if v.GetId() == reject.GetOperatorId() {
			signatoryIndex = i
			break
		}
	}
	if signatoryIndex < 0 {
		return nil, cty.ErrNotInSignatories
	}
	if contract.GetSignatories()[signatoryIndex].GetSignStatus() != cty.SignStatus_SS_Unsigned {
		return nil, cty.ErrSignStatus
	}
	contract.Signatories[signatoryIndex].SignStatus = cty.SignStatus_SS_Rejected
	contract.Signatories[signatoryIndex].RejectTime = reject.GetOperateTime()
	contract.Status = cty.ContractStatus_CS_Disused
	contract.UpdateTime = reject.GetOperateTime()

	c.contract = *contract
	kv = append(kv, c.save(action.db)...)
	logs = append(logs, &types.ReceiptLog{Ty: cty.TyLogRejectContract, Log: types.Encode(contract)})

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

//func (action *Action) userRegister(register *cty.ContractRegister) (*types.Receipt, error) {
//	var logs []*types.ReceiptLog
//	var kv []*types.KeyValue
//	c := newDB()
//
//	//数据格式验证
//	if err := register.Validate(); err != nil {
//		return nil, cty.ErrFormat
//	}
//
//	//判断用户是否已存在
//	_, err := action.getUserById(register.GetId())
//	if err != types.ErrNotFound {
//		if err != nil {
//			return nil, err
//		}
//		return nil, cty.ErrUserExists
//	}
//	user := &cty.User{
//		Id:     register.GetId(),
//		Name:   register.GetName(),
//		PubKey: register.GetPubKey(),
//	}
//
//	c.user = *user
//	kv = append(kv, c.save(action.db)...)
//
//	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
//}

func isContractEffective(contract *cty.Contract) bool {
	if contract.GetCreator().GetSignStatus() != cty.SignStatus_SS_Signed {
		return false
	}
	for _, v := range contract.GetSignatories() {
		if v.GetSignStatus() != cty.SignStatus_SS_Signed {
			return false
		}
	}
	return true
}

func (action *Action) getContractById(contractId string) (*cty.Contract, error) {
	value, err := action.db.Get(calcContractKey(contractId))
	if err != nil {
		return nil, err
	}

	var info cty.Contract
	if err = types.Decode(value, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (action *Action) getUserById(userId string) ([]byte, error) {
	return action.db.Get(calcUserKey(userId))
}

//判断是否重复
func isDuplicate(args ...interface{}) bool {
	for i, v1 := range args {
		for _, v2 := range args[i+1:] {
			if v1 == v2 {
				return true
			}
		}
	}
	return false
}
