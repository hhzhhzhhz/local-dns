package constant

type ApiResponse struct {
	Code int
	Info string
}

func NewResponse(code int, info string) ApiResponse {
	return ApiResponse{code, info}
}

var (
	// ErrorParams params 1-50
	ErrorParams = NewResponse(1, "invalid params")
	ErrorReadBody = NewResponse(2, "read failed")

	// Common 51- 100
	ErrorDbOp = NewResponse(51, "db op failed")
)
