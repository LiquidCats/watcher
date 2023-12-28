package bitcoin

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

var hash = "00000000000000000002c1d4c923fdf69aa283f242df191e0159cb1143958155"
var previous = "000000000000000000013bbbd5ff9f606f8bf543512e286ee17544ca0ce0a32c"
var height = 822_602

func (r *rpcCaller) CallFor(ctx context.Context, dest any, method string, params ...any) error {
	if method == "getblockheader" {
		header := dest.(*blockHeader)

		header.Hash = hash
		header.PreviousBlockHash = previous
		header.Number = height
	}

	args := r.Called(ctx, dest, method, params)
	return args.Error(0)
}

func TestRpcRepository_GetLatestBlock(t *testing.T) {
	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	ctx := context.Background()

	mockedClient.On("CallFor", ctx, mock.Anything, "getbestblockhash", mock.Anything).Once().Return(nil)
	mockedClient.On("CallFor", ctx, mock.Anything, "getblockheader", mock.Anything).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetLatestBlock(ctx)

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)

}

func TestRpcRepository_GetBlockByHash(t *testing.T) {
	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	ctx := context.Background()

	mockedClient.On("CallFor", ctx, mock.Anything, "getblockheader", mock.Anything).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetBlockByHash(ctx, entity.BlockHash(hash))

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)
}

func TestRpcRepository_GetBlockByHeight(t *testing.T) {
	mockedClient := new(rpcCaller)
	defer mockedClient.AssertExpectations(t)

	ctx := context.Background()

	mockedClient.On("CallFor", ctx, mock.Anything, "getblockstats", mock.Anything).Once().Return(nil)
	mockedClient.On("CallFor", ctx, mock.Anything, "getblockheader", mock.Anything).Once().Return(nil)

	repo := NewRpcRepository(mockedClient)

	block, err := repo.GetBlockByHeight(ctx, entity.BlockHeight(height))

	require.NoError(t, err)

	require.Equal(t, entity.BlockHeight(height), block.Height)
	require.Equal(t, entity.BlockHash(hash), block.Hash)
	require.Equal(t, entity.BlockHash(previous), block.Previous)
}
