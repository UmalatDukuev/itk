package errs

import "errors"

var (
	ErrInvalidOperationType = errors.New("invalid operation type")
	ErrInsufficientFunds    = errors.New("insufficient funds")
	ErrWalletNotFound       = errors.New("wallet not found")
)
