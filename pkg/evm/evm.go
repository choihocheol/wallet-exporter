package evm

import (
	"context"
	"encoding/json"
	"fmt"
	"main/pkg/config"
	"main/pkg/http"
	"main/pkg/types"
	"math/big"
	"strings"
	"sync"

	"cosmossdk.io/math"
	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

type RPC struct {
	Client *http.Client
	URL    string
	Denoms []config.DenomInfo
	Logger zerolog.Logger
	Tracer trace.Tracer

	LastQueryHeight map[string]int64
	Mutex           sync.Mutex
}

func NewRPC(chain config.Chain, logger zerolog.Logger, tracer trace.Tracer) *RPC {
	return &RPC{
		Client:          http.NewClient(logger, chain.Name, tracer),
		URL:             chain.QueryEndpoint,
		Denoms:          chain.Denoms,
		Logger:          logger.With().Str("component", "rpc").Logger(),
		LastQueryHeight: make(map[string]int64),
		Tracer:          tracer,
	}
}

func (rpc *RPC) GetWalletBalances(
	address string,
	ctx context.Context,
) (*types.BalanceResponse, types.QueryInfo, error) {
	req := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_getBalance",
		Params: []any{
			address,
			"latest",
		},
		ID: 1,
	}
	var resp RPCResponse

	queryInfo, _, err := rpc.Client.Post(rpc.URL, req, &resp, accept200, ctx)
	if err != nil {
		return nil, queryInfo, err
	}

	if resp.Error != nil {
		queryInfo.Success = false
		return nil, queryInfo, fmt.Errorf(
			"rpc error (%d): %s",
			resp.Error.Code,
			resp.Error.Message,
		)
	}

	var balanceHex string
	if err := json.Unmarshal(resp.Result, &balanceHex); err != nil {
		queryInfo.Success = false
		return nil, queryInfo, fmt.Errorf("failed to decode balance: %w", err)
	}

	balance, ok := new(big.Int).SetString(strings.TrimPrefix(balanceHex, "0x"), 16)
	if !ok {
		queryInfo.Success = false
		return nil, queryInfo, fmt.Errorf("invalid balance: %s", balanceHex)
	}

	amount := math.LegacyNewDecFromBigInt(balance)

	return &types.BalanceResponse{
		Balances: types.Balances{
			{
				Denom:  rpc.Denoms[0].Denom,
				Amount: amount,
			},
		},
	}, queryInfo, nil
}
