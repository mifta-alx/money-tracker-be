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
	"invalid_credentials":  "Invalid email or password",
	"user_not_found":       "User account not found.",
	"unauthorized":         "You are not authorized to access this resource.",
	"invalid_token":        "Invalid or expired token, please login again.",
	"auth_header_required": "Authorization header is required.",

	// OAuth / Google Login Errors
	"failed_create_oauth":        "Failed to create account via Google.",
	"failed_link_google_account": "Failed to link your Google account to the existing profile.",
	"google_auth_failed":         "Google authentication failed.",
	"token_required":             "Token is required",
	"invalid_google_token":       "Invalid token, please try again later",

	// General Errors
	"internal_server_error": "An unexpected error occurred on our server",
	"validation_error":      "Validation failed",
	"malformed_request":     "Malformed request",

	// Accounts Errors
	"account_limit_reached": "You have reached the maximum number of accounts (10)",
	"account_not_found":     "Account not found",

	// Category Errors
	"category_not_found": "Category not found",
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
