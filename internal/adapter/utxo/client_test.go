package utxo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/utxo"
	"github.com/LiquidCats/watcher/v2/internal/app/kernel/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	rawMempoolResponse  = `{"jsonrpc":"2.0","result":["c27b0cd13a4159143b43502f47d18a308e3f995902d123912b470d0bef636738","78f3731b3a5a83bbaa4062cdad1bcbffa98951138c98c09cd35f278dbcfbaace","e1bf96fe40df1366dae7cd56f0dc07947f64cc09676d52b576fb4d3c64d94dc8","6e23d2635c9f9c862270ca7c850bc86158524f2adf9c0b5eb0ece6482d820fbd","7ff55d8ea89c97baefe624573eed5c1d7caa63ca94a01d5519a20302f2c505de","983db2078d0fe2ae0578e2fb6169f43ac7f815847dc5be7e4bf6f192b655896b"]},"id":"getblock.io"}`
	blockResponse       = `{"jsonrpc":"2.0","result":{"hash":"00000000000000000001246cadf2834cf70fe92404e37c14e071f4b7da61993d","confirmations":3,"height":883189,"version":560087040,"versionHex":"21624000","merkleroot":"93960bafc07fe95739ca835f5d1b30b78634d9859ea78ae667f94a01c1f289a7","time":1739204667,"mediantime":1739197922,"nonce":528725547,"bits":"17027726","difficulty":114167270716407.6,"chainwork":"0000000000000000000000000000000000000000ad4e5f6f4e49d96c1e4a2cdc","nTx":3125,"previousblockhash":"00000000000000000001778bcbd65816f9d0944345448f4fa19f17946ab59078","nextblockhash":"0000000000000000000189e4159043cf8e807c7d181cdcc90b4b7d0db4ed522e","strippedsize":759467,"size":1714974,"weight":3993375,"tx":["a3667aba97984ad2723d955a1817c35563f0c5d1af2f96f005c6a00ad612a5cc","18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e","70b5d862dcbe448701069cc323f112f16cf291dcb16c47f56b1aea61640bcd4e","983db2078d0fe2ae0578e2fb6169f43ac7f815847dc5be7e4bf6f192b655896b"]},"id":"getblock.io"}`
	transactionResponse = `{"jsonrpc":"2.0","result":{"txid":"18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e","hash":"925f67547a7a5b6cab181e866a591ab632468da00315cf02243eb299b954fc19","version":2,"size":234,"vsize":153,"weight":609,"locktime":0,"vin":[{"txid":"d5aa0a84cc4cdcfbb2ffb175abf638a45468fbf8eb9a2d1aa380e6e41ac1030b","vout":4,"scriptSig":{"asm":"","hex":""},"txinwitness":["304402202b88a9317decf5e219ed44f6947ca28fec4e2ad378496d5b5d6196a7c253518d022017503642ecd2b2b451e68674d1ea2b2006432809cd32d613a59e21221b49b35801","033b3b26324525c22e39fe7942bcac87976804b776c841bbf51f05d22e75a269c7"],"sequence":4294967293}],"vout":[{"value":0.00295078,"n":0,"scriptPubKey":{"asm":"0 3753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541af","desc":"addr(bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z)#9sh3632j","hex":"00203753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541af","address":"bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z","type":"witness_v0_scripthash"}},{"value":0.02143401,"n":1,"scriptPubKey":{"asm":"0 b84ef1e7e398d56fa88a356df3139ff997114f04","desc":"addr(bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c)#3cnwxrlh","hex":"0014b84ef1e7e398d56fa88a356df3139ff997114f04","address":"bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c","type":"witness_v0_keyhash"}}],"hex":"020000000001010b03c11ae4e680a31a2d9aebf8fb6854a438f6ab75b1ffb2fbdc4ccc840aaad50400000000fdffffff02a6800400000000002200203753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541afa9b4200000000000160014b84ef1e7e398d56fa88a356df3139ff997114f040247304402202b88a9317decf5e219ed44f6947ca28fec4e2ad378496d5b5d6196a7c253518d022017503642ecd2b2b451e68674d1ea2b2006432809cd32d613a59e21221b49b3580121033b3b26324525c22e39fe7942bcac87976804b776c841bbf51f05d22e75a269c700000000","blockhash":"00000000000000000001246cadf2834cf70fe92404e37c14e071f4b7da61993d","confirmations":713,"time":1739204667,"blocktime":1739204667},"id":"getblock.io"}`
)

func TestClient_GetMempool(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(rawMempoolResponse))
	}))
	defer api.Close()

	client := utxo.NewClient[entities.TxID](configs.Utxo{NodeUrl: api.URL})

	ctx := context.Background()

	result, err := client.GetMempool(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	assert.Len(t, result, 6)
}

func TestClient_GetBlockByHash(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(blockResponse))
	}))
	defer api.Close()

	client := utxo.NewClient[entities.TxID](configs.Utxo{NodeUrl: api.URL})

	ctx := context.Background()

	result, err := client.GetBlockByHash(ctx, "00000000000000000001246cadf2834cf70fe92404e37c14e071f4b7da61993d")
	require.NoError(t, err)

	assert.Equal(t, entities.BlockHash("00000000000000000001246cadf2834cf70fe92404e37c14e071f4b7da61993d"), result.Hash)
	assert.NotEmpty(t, result.Tx)
	assert.Len(t, result.Tx, 4)
}

func TestClient_GetTransactionByTxId(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(transactionResponse))
	}))
	defer api.Close()

	client := utxo.NewClient[entities.TxID](configs.Utxo{NodeUrl: api.URL})

	ctx := context.Background()

	result, err := client.GetTransactionByTxId(ctx, "18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e")
	require.NoError(t, err)

	assert.Equal(t, entities.TxID("18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e"), result.Txid)

	assert.NotEmpty(t, result.Vin)
	assert.Len(t, result.Vin, 1)

	assert.Equal(t, entities.TxID("d5aa0a84cc4cdcfbb2ffb175abf638a45468fbf8eb9a2d1aa380e6e41ac1030b"), result.Vin[0].Txid)

	assert.NotEmpty(t, result.Vout)
	assert.Len(t, result.Vout, 2)

	assert.Equal(t, entities.Address("bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z"), result.Vout[0].ScriptPubKey.Address)
	assert.Equal(t, entities.Address("bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c"), result.Vout[1].ScriptPubKey.Address)
}
