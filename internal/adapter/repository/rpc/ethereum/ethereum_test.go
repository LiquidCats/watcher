package ethereum

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"watcher/internal/app/domain/entity"
)

type rpcCaller struct {
	mock.Mock
}

var hash = "0xfb1c2738f5571163fdc82657eb71895ffaf4289a57c608585442b52084fb709b"
var previous = "0xfdd20ea215322b858bdd0e328c5dcf55de53e54d2217b9da9199997de92069e9"
var height = 18_850_374

func (r *rpcCaller) CallFor(ctx context.Context, dest any, method string, params ...any) error {
	b := dest.(*block)
	b.Height = toHex(entity.BlockHeight(height))
	b.Hash = hash
	b.ParentHash = previous

	args := r.Called(ctx, dest, method, params)
	return args.Error(0)
}

func TestRpcRepository_GetBlockByHeight(t *testing.T) {
	ctx := context.Background()

	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	mockedClient.On("CallFor", ctx, mock.Anything, "eth_getBlockByNumber", []any{toHex(entity.BlockHeight(height)), false}).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetBlockByHeight(ctx, entity.BlockHeight(height))

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)
}

func TestRpcRepository_GetBlockByHash(t *testing.T) {
	ctx := context.Background()

	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	mockedClient.On("CallFor", ctx, mock.Anything, "eth_getBlockByHash", []any{hash, false}).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetBlockByHash(ctx, entity.BlockHash(hash))

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)
}

func TestRpcRepository_GetLatestBlock(t *testing.T) {
	ctx := context.Background()

	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	mockedClient.On("CallFor", ctx, mock.Anything, "eth_getBlockByNumber", []any{"latest", false}).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetLatestBlock(ctx)

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)
}
