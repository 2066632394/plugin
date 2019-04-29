package types

import "errors"

var (
	ErrContractExists   = errors.New("ErrContractExists")
	ErrSignStatus       = errors.New("ErrSignStatus")
	ErrContractStatus   = errors.New("ErrContractStatus")
	ErrFormat           = errors.New("ErrFormat")
	ErrPubKey           = errors.New("ErrPubKey")
	ErrHash             = errors.New("ErrHash")
	ErrPermissionDenied = errors.New("ErrPermissionDenied")
	ErrOperateTime      = errors.New("ErrOperateTime")
	ErrSignedHash       = errors.New("ErrSignedHash")
	ErrDuplicateUserId  = errors.New("ErrDuplicateUserId")
	ErrNotInSignatories = errors.New("ErrNotInSignatories")
	ErrUserExists       = errors.New("ErrUserExists")
	ErrSameHash         = errors.New("ErrSameHash")
)
