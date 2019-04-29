package types

import (
	"github.com/33cn/chain33/types"
	"reflect"
)

var (
	ContractX      = "contract"
	ExecerContract = []byte(ContractX)
	//ManagerPubKey  = "4be3ca87cab28f156516a65cb91acf1d33020f3c4214005dd95d44ff872062d3"
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerContract)
	types.RegistorExecutor(ContractX, NewType())
}

type ContractType struct {
	types.ExecTypeBase
}

func NewType() *ContractType {
	c := &ContractType{}
	c.SetChild(c)
	return c
}

func (c *ContractType) GetPayload() types.Message {
	return &ContractAction{}
}

func (c *ContractType) GetName() string {
	return ContractX
}

func (c *ContractType) GetLogMap() map[int64]*types.LogInfo {
	return map[int64]*types.LogInfo{
		TyLogCreateContract: {reflect.TypeOf(Contract{}), "LogCreateContract"},
		TyLogCancelContract: {reflect.TypeOf(Contract{}), "LogCancelContract"},
		TyLogModifyContract: {reflect.TypeOf(Contract{}), "LogModifyContract"},
		TyLogSignContract:   {reflect.TypeOf(Contract{}), "LogSignContract"},
		TyLogRejectContract: {reflect.TypeOf(Contract{}), "LogRejectContract"},
	}
}

func (c *ContractType) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Create": ContractActionCreate,
		"Cancel": ContractActionCancel,
		"Modify": ContractActionModify,
		"Sign":   ContractActionSign,
		"Reject": ContractActionReject,
		//"Register": ContractActionRegister,
	}
}

func (c *ContractType) ActionName(tx *types.Transaction) string {
	var action ContractAction
	err := types.Decode(tx.Payload, &action)
	if err != nil {
		return "unknown"
	}
	switch {
	case action.Ty == ContractActionCreate && action.GetCreate() != nil:
		return "Create"
	case action.Ty == ContractActionCancel && action.GetCancel() != nil:
		return "Cancel"
	case action.Ty == ContractActionModify && action.GetModify() != nil:
		return "Modify"
	case action.Ty == ContractActionSign && action.GetSign() != nil:
		return "Sign"
	case action.Ty == ContractActionReject && action.GetReject() != nil:
		return "Reject"
		//case action.Ty == ContractActionRegister && action.GetRegister() != nil:
		//	return "Register"
	}

	return "unknown"
}

func (c *ContractType) Amount(tx *types.Transaction) (int64, error) {
	return 0, nil
}
