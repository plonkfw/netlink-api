package types

// APIResponse standardization
type APIResponse struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    []byte `json:"data"`
}
