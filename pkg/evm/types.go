package evm

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

type RPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func accept200(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", res.StatusCode)
	}
	return nil
}
