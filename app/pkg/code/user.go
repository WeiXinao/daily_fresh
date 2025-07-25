package code

//go:generate codegen -type=int -doc -output ./error_code_generated.md
const (
  // ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists 

	// ErrLoginFailed - 500: Login failed.
	ErrLoginFailed

	// ErrSmsSend - 500: Send sms error.
	ErrSmsSend

	// ErrCodeNotExist - 400: Sms code not exist.
	ErrCodeNotExist

	// ErrCodeIncorrent - 400: Sms code incorrent.
	ErrCodeIncorrent
)