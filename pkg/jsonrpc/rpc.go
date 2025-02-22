package jsonrpc

import "time"

type RPCResponse[D any] struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  D         `json:"result"`
	Error   *RPCError `json:"error,omitempty"`
	ID      any       `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RPCRequest[P any] struct {
	Method  string `json:"method"`
	Params  P      `json:"params,omitempty"`
	ID      int64  `json:"id"`
	JSONRPC string `json:"jsonrpc"`
}

func createRequest[P any](method string, params P) *RPCRequest[P] {
	return &RPCRequest[P]{
		ID:      time.Now().UnixMilli(),
		Method:  method,
		JSONRPC: "2.0",
		Params:  params,
	}
}
