package bitcoin

import (
	"context"
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
	var hash entity.BlockHash
	if err := r.client.CallFor(ctx, &hash, "getbestblockhash"); err != nil {
		return nil, err
	}

	return r.GetBlockByHash(ctx, hash)
}

func (r *RpcRepository) GetBlockByHash(ctx context.Context, hash entity.BlockHash) (*entity.Block, error) {
	var header blockHeader
	if err := r.client.CallFor(ctx, &header, "getblockheader", string(hash)); err != nil {
		return nil, err
	}

	return header.toEntity(), nil
}

func (r *RpcRepository) GetBlockByHeight(ctx context.Context, height entity.BlockHeight) (*entity.Block, error) {
	var stats blockStats
	if err := r.client.CallFor(ctx, &stats, "getblockstats", int(height)); err != nil {
		return nil, err
	}

	return r.GetBlockByHash(ctx, entity.BlockHash(stats.Hash))
}

func (r *RpcRepository) GetBlockTransactions(_ context.Context, _ entity.BlockHash) ([]entity.Transaction, error) {
	return nil, errors.ErrNotImplemented
}
