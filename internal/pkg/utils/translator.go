package utils

var errorMessages = map[string]string{
	// Registration Errors
	"email_exists":            "Email is already registered, please use another one.",
	"password_too_short":      "Password must be at least 8 characters long.",
	"password_contains_space": "Password must not contain any spaces.",
	"password_no_upper":       "Password must contain at least one uppercase letter.",
	"password_no_lower":       "Password must contain at least one lowercase letter.",
	"password_no_number":      "Password must contain at least one number.",
	"process_password_failed": "Failed to process security credentials.",
	"create_account_failed":   "Could not create account, please try again later.",

	// Login & Auth Errors
	"invalid_credentials": "Invalid email or password.",
	"user_not_found":      "User account not found.",
	"unauthorized":        "You are not authorized to access this resource.",

	// General Errors
	"internal_server_error": "An unexpected error occurred on our server.",
}

func TranslateError(err error) string {
	if err == nil {
		return ""
	}

	msg, ok := errorMessages[err.Error()]
	if !ok {
		return "An unexpected error occurred."
	}

	return msg
}
