package services

import "errors"

var (
	ErrAccountNotFound         = errors.New("account_not_found")
	ErrCategoryNotFound        = errors.New("category_not_found")
	ErrTransactionNotFound     = errors.New("transaction_not_found")
	ErrUserNotFound            = errors.New("user_not_found")
	ErrInternal                = errors.New("internal_server_error")
	ErrAccountLimitReached     = errors.New("account_limit_reached")
	ErrMissingRequiredFields   = errors.New("missing_required_fields")
	ErrBalanceCannotBeNegative = errors.New("balance_cannot_be_negative")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrValidation              = errors.New("validation_error")
	ErrMalformedRequest        = errors.New("malformed_request")
	ErrEmailExist              = errors.New("email_exists")
	ErrPasswordTooShort        = errors.New("password_too_short")
	ErrPasswordContainsSpace   = errors.New("password_contains_space")
	ErrPasswordNoNumber        = errors.New("password_no_number")
	ErrPasswordNoUpper         = errors.New("password_no_upper")
	ErrPasswordNoLower         = errors.New("password_no_lower")
	ErrInvalidCredentials      = errors.New("invalid_credentials")
	ErrFailedCreateOAuth       = errors.New("failed_create_oauth")
	ErrFailedLinkGoogle        = errors.New("failed_link_google_account")
	ErrTokenRequired           = errors.New("token_required")
	ErrInvalidGoogleToken      = errors.New("invalid_google_token")
)
