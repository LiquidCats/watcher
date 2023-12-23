package ethereum

import (
	"context"
	"fmt"
	"strconv"
	"watcher/internal/app/domain/entity"
	"watcher/internal/app/domain/errors"
	"watcher/internal/port"
)

type RpcRepository struct {
	client port.RpcCaller
}

func NewRpcRepository(client port.RpcCaller) *RpcRepository {
	return &RpcRepository{
		client: client,
	}
}

func (r *RpcRepository) GetLatestBlock(ctx context.Context) (*entity.Block, error) {
	var latest block
	if err := r.client.CallFor(ctx, &latest, "eth_getBlockByNumber", "latest", false); err != nil {
		return nil, err
	}

	return latest.toEntity(), nil
}

func (r *RpcRepository) GetBlockByHash(ctx context.Context, hash entity.BlockHash) (*entity.Block, error) {
	var b block
	if err := r.client.CallFor(ctx, &b, "eth_getBlockByHash", string(hash), false); err != nil {
		return nil, err
	}

	return b.toEntity(), nil
}

func (r *RpcRepository) GetBlockByHeight(ctx context.Context, height entity.BlockHeight) (*entity.Block, error) {
	var b block
	if err := r.client.CallFor(ctx, &b, "eth_getBlockByNumber", toHex(height), false); err != nil {
		return nil, err
	}

	return b.toEntity(), nil
}

func (r *RpcRepository) GetBlockTransactions(_ context.Context, _ entity.BlockHash) ([]entity.Transaction, error) {
	return nil, errors.ErrNotImplemented
}

func toHex(height entity.BlockHeight) string {
	hex := strconv.FormatInt(int64(height), 16)
	return fmt.Sprint("0x", hex)
}
