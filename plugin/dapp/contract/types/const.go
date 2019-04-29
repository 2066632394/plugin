package types

const (
	ContractActionCreate = iota + 1
	ContractActionCancel
	ContractActionModify
	ContractActionSign
	ContractActionReject
	//ContractActionRegister
)

const (
	TyLogCreateContract = iota + 1100
	TyLogCancelContract
	TyLogModifyContract
	TyLogSignContract
	TyLogRejectContract
)
