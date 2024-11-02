package service_errors

const (
	//Token
	UnexpectedError = "Expected error"
	ClaimsNotFound  = "Claims not found"
	TokenRequired   = "Token is required"
	TokenExpired    = "Token is expired"
	TokenInvalid    = "Token is invalid"

	// OTP
	OtpExists  = "Otp Exists"
	OtpUsed    = "Otp Used"
	OtpInvalid = "Otp Invalid"

	//User
	EmailExists    = "Email exists"
	UsernameExists = "Username exists"

	// DB
	RecordNotFound = "record not found"
)
