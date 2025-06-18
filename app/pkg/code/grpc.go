package code

//go:generate codegen -type=int -doc -output ./error_code_generated.md
const (
  // ErrGrpcConn - 500: Connect to grpc error.
	ErrGrpcConn int = iota + 100601
)