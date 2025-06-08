package code

//go:generate codegen -type=int -doc -output ./error_code_generated.md
const (
  // ErrUserNotFound - 404: User not found.
	ErrUserNotFound int = iota + 100401

	// ErrUserAlreadyExists - 400: User already exists.
	ErrUserAlreadyExists 
)