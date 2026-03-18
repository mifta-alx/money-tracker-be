package services

import "errors"

var (
	ErrAccountNotFound         = errors.New("account_not_found")
	ErrInternal                = errors.New("internal_server_error")
	ErrAccountLimitReached     = errors.New("account_limit_reached")
	ErrMissingRequiredFields   = errors.New("missing_required_fields")
	ErrBalanceCannotBeNegative = errors.New("balance_cannot_be_negative")
)
