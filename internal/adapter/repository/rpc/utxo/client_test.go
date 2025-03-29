package utxo_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LiquidCats/jsonrpc"
	"github.com/LiquidCats/watcher/v2/configs"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo"
	"github.com/LiquidCats/watcher/v2/internal/adapter/repository/rpc/utxo/data"
	"github.com/LiquidCats/watcher/v2/internal/app/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	rawMempoolResponse  = `{"jsonrpc":"2.0","result":["c27b0cd13a4159143b43502f47d18a308e3f995902d123912b470d0bef636738","78f3731b3a5a83bbaa4062cdad1bcbffa98951138c98c09cd35f278dbcfbaace","e1bf96fe40df1366dae7cd56f0dc07947f64cc09676d52b576fb4d3c64d94dc8","6e23d2635c9f9c862270ca7c850bc86158524f2adf9c0b5eb0ece6482d820fbd","7ff55d8ea89c97baefe624573eed5c1d7caa63ca94a01d5519a20302f2c505de","983db2078d0fe2ae0578e2fb6169f43ac7f815847dc5be7e4bf6f192b655896b"]},"id":"getblock.io"}`
	blockResponse       = `{"jsonrpc":"2.0","result":{"hash":"000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506","confirmations":785765,"height":100000,"version":1,"versionHex":"00000001","merkleroot":"f3e94742aca4b5ef85488dc37c06c3282295ffec960994b2c0d5ac2a25a95766","time":1293623863,"mediantime":1293622620,"nonce":274148111,"bits":"1b04864c","difficulty":14484.1623612254,"chainwork":"0000000000000000000000000000000000000000000000000644cb7f5234089e","nTx":4,"previousblockhash":"000000000002d01c1fccc21636b607dfd930d31d01c3a62104612a1719011250","nextblockhash":"00000000000080b66c911bd5ba14a74260057311eaeb1982802f7010f1a9f090","strippedsize":957,"size":957,"weight":3828,"tx":[{"txid":"8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87","hash":"8c14f0db3df150123e6f3dbbf30f8b955a8249b62ac1d1ff16284aefa3d06d87","version":1,"size":135,"vsize":135,"weight":540,"locktime":0,"vin":[{"coinbase":"044c86041b020602","sequence":4294967295}],"vout":[{"value":50.00000000,"n":0,"scriptPubKey":{"asm":"041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84 OP_CHECKSIG","desc":"pk(041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84)#40d2kraw","hex":"41041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac","type":"pubkey"}}],"hex":"01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff08044c86041b020602ffffffff0100f2052a010000004341041b0e8c2567c12536aa13357b79a073dc4444acb83c4ec7a0e2f99dd7457516c5817242da796924ca4e99947d087fedf9ce467cb9f7c6287078f801df276fdf84ac00000000"},{"txid":"fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4","hash":"fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4","version":1,"size":259,"vsize":259,"weight":1036,"locktime":0,"vin":[{"txid":"87a157f3fd88ac7907c05fc55e271dc4acdc5605d187d646604ca8c0e9382e03","vout":0,"scriptSig":{"asm":"3046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748[ALL] 04f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3","hex":"493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3"},"sequence":4294967295}],"vout":[{"value":5.56000000,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 c398efa9c392ba6013c5e04ee729755ef7f58b32 OP_EQUALVERIFY OP_CHECKSIG","desc":"addr(1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn)#70tf4ktg","hex":"76a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac","address":"1JqDybm2nWTENrHvMyafbSXXtTk5Uv5QAn","type":"pubkeyhash"}},{"value":44.44000000,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 948c765a6914d43f2a7ac177da2c2f6b52de3d7c OP_EQUALVERIFY OP_CHECKSIG","desc":"addr(1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx)#t40fmm8a","hex":"76a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac","address":"1EYTGtG4LnFfiMvjJdsU7GMGCQvsRSjYhx","type":"pubkeyhash"}}],"fee":0.00000000,"hex":"0100000001032e38e9c0a84c6046d687d10556dcacc41d275ec55fc00779ac88fdf357a187000000008c493046022100c352d3dd993a981beba4a63ad15c209275ca9470abfcd57da93b58e4eb5dce82022100840792bc1f456062819f15d33ee7055cf7b5ee1af1ebcc6028d9cdb1c3af7748014104f46db5e9d61a9dc27b8d64ad23e7383a4e6ca164593c2527c038c0857eb67ee8e825dca65046b82c9331586c82e0fd1f633f25f87c161bc6f8a630121df2b3d3ffffffff0200e32321000000001976a914c398efa9c392ba6013c5e04ee729755ef7f58b3288ac000fe208010000001976a914948c765a6914d43f2a7ac177da2c2f6b52de3d7c88ac00000000"},{"txid":"6359f0868171b1d194cbee1af2f16ea598ae8fad666d9b012c8ed2b79a236ec4","hash":"6359f0868171b1d194cbee1af2f16ea598ae8fad666d9b012c8ed2b79a236ec4","version":1,"size":257,"vsize":257,"weight":1028,"locktime":0,"vin":[{"txid":"cf4e2978d0611ce46592e02d7e7daf8627a316ab69759a9f3df109a7f2bf3ec3","vout":1,"scriptSig":{"asm":"30440220032d30df5ee6f57fa46cddb5eb8d0d9fe8de6b342d27942ae90a3231e0ba333e02203deee8060fdc70230a7f5b4ad7d7bc3e628cbe219a886b84269eaeb81e26b4fe[ALL] 04ae31c31bf91278d99b8377a35bbce5b27d9fff15456839e919453fc7b3f721f0ba403ff96c9deeb680e5fd341c0fc3a7b90da4631ee39560639db462e9cb850f","hex":"4730440220032d30df5ee6f57fa46cddb5eb8d0d9fe8de6b342d27942ae90a3231e0ba333e02203deee8060fdc70230a7f5b4ad7d7bc3e628cbe219a886b84269eaeb81e26b4fe014104ae31c31bf91278d99b8377a35bbce5b27d9fff15456839e919453fc7b3f721f0ba403ff96c9deeb680e5fd341c0fc3a7b90da4631ee39560639db462e9cb850f"},"sequence":4294967295}],"vout":[{"value":0.01000000,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 b0dcbf97eabf4404e31d952477ce822dadbe7e10 OP_EQUALVERIFY OP_CHECKSIG","desc":"addr(1H8ANdafjpqYntniT3Ddxh4xPBMCSz33pj)#ejxdpym6","hex":"76a914b0dcbf97eabf4404e31d952477ce822dadbe7e1088ac","address":"1H8ANdafjpqYntniT3Ddxh4xPBMCSz33pj","type":"pubkeyhash"}},{"value":2.99000000,"n":1,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 6b1281eec25ab4e1e0793ff4e08ab1abb3409cd9 OP_EQUALVERIFY OP_CHECKSIG","desc":"addr(1Am9UTGfdnxabvcywYG2hvzr6qK8T3oUZT)#te23gwt4","hex":"76a9146b1281eec25ab4e1e0793ff4e08ab1abb3409cd988ac","address":"1Am9UTGfdnxabvcywYG2hvzr6qK8T3oUZT","type":"pubkeyhash"}}],"fee":0.00000000,"hex":"0100000001c33ebff2a709f13d9f9a7569ab16a32786af7d7e2de09265e41c61d078294ecf010000008a4730440220032d30df5ee6f57fa46cddb5eb8d0d9fe8de6b342d27942ae90a3231e0ba333e02203deee8060fdc70230a7f5b4ad7d7bc3e628cbe219a886b84269eaeb81e26b4fe014104ae31c31bf91278d99b8377a35bbce5b27d9fff15456839e919453fc7b3f721f0ba403ff96c9deeb680e5fd341c0fc3a7b90da4631ee39560639db462e9cb850fffffffff0240420f00000000001976a914b0dcbf97eabf4404e31d952477ce822dadbe7e1088acc060d211000000001976a9146b1281eec25ab4e1e0793ff4e08ab1abb3409cd988ac00000000"},{"txid":"e9a66845e05d5abc0ad04ec80f774a7e585c6e8db975962d069a522137b80c1d","hash":"e9a66845e05d5abc0ad04ec80f774a7e585c6e8db975962d069a522137b80c1d","version":1,"size":225,"vsize":225,"weight":900,"locktime":0,"vin":[{"txid":"f4515fed3dc4a19b90a317b9840c243bac26114cf637522373a7d486b372600b","vout":0,"scriptSig":{"asm":"3046022100bb1ad26df930a51cce110cf44f7a48c3c561fd977500b1ae5d6b6fd13d0b3f4a022100c5b42951acedff14abba2736fd574bdb465f3e6f8da12e2c5303954aca7f78f3[ALL] 04a7135bfe824c97ecc01ec7d7e336185c81e2aa2c41ab175407c09484ce9694b44953fcb751206564a9c24dd094d42fdbfdd5aad3e063ce6af4cfaaea4ea14fbb","hex":"493046022100bb1ad26df930a51cce110cf44f7a48c3c561fd977500b1ae5d6b6fd13d0b3f4a022100c5b42951acedff14abba2736fd574bdb465f3e6f8da12e2c5303954aca7f78f3014104a7135bfe824c97ecc01ec7d7e336185c81e2aa2c41ab175407c09484ce9694b44953fcb751206564a9c24dd094d42fdbfdd5aad3e063ce6af4cfaaea4ea14fbb"},"sequence":4294967295}],"vout":[{"value":0.01000000,"n":0,"scriptPubKey":{"asm":"OP_DUP OP_HASH160 39aa3d569e06a1d7926dc4be1193c99bf2eb9ee0 OP_EQUALVERIFY OP_CHECKSIG","desc":"addr(16FuTPaeRSPVxxCnwQmdyx2PQWxX6HWzhQ)#r33agg0y","hex":"76a91439aa3d569e06a1d7926dc4be1193c99bf2eb9ee088ac","address":"16FuTPaeRSPVxxCnwQmdyx2PQWxX6HWzhQ","type":"pubkeyhash"}}],"fee":0.00000000,"hex":"01000000010b6072b386d4a773235237f64c1126ac3b240c84b917a3909ba1c43ded5f51f4000000008c493046022100bb1ad26df930a51cce110cf44f7a48c3c561fd977500b1ae5d6b6fd13d0b3f4a022100c5b42951acedff14abba2736fd574bdb465f3e6f8da12e2c5303954aca7f78f3014104a7135bfe824c97ecc01ec7d7e336185c81e2aa2c41ab175407c09484ce9694b44953fcb751206564a9c24dd094d42fdbfdd5aad3e063ce6af4cfaaea4ea14fbbffffffff0140420f00000000001976a91439aa3d569e06a1d7926dc4be1193c99bf2eb9ee088ac00000000"}]},"id":"getblock.io"}`
	transactionResponse = `{"jsonrpc":"2.0","result":{"txid":"18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e","hash":"925f67547a7a5b6cab181e866a591ab632468da00315cf02243eb299b954fc19","version":2,"size":234,"vsize":153,"weight":609,"locktime":0,"vin":[{"txid":"d5aa0a84cc4cdcfbb2ffb175abf638a45468fbf8eb9a2d1aa380e6e41ac1030b","vout":4,"scriptSig":{"asm":"","hex":""},"txinwitness":["304402202b88a9317decf5e219ed44f6947ca28fec4e2ad378496d5b5d6196a7c253518d022017503642ecd2b2b451e68674d1ea2b2006432809cd32d613a59e21221b49b35801","033b3b26324525c22e39fe7942bcac87976804b776c841bbf51f05d22e75a269c7"],"sequence":4294967293}],"vout":[{"value":0.00295078,"n":0,"scriptPubKey":{"asm":"0 3753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541af","desc":"addr(bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z)#9sh3632j","hex":"00203753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541af","address":"bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z","type":"witness_v0_scripthash"}},{"value":0.02143401,"n":1,"scriptPubKey":{"asm":"0 b84ef1e7e398d56fa88a356df3139ff997114f04","desc":"addr(bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c)#3cnwxrlh","hex":"0014b84ef1e7e398d56fa88a356df3139ff997114f04","address":"bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c","type":"witness_v0_keyhash"}}],"hex":"020000000001010b03c11ae4e680a31a2d9aebf8fb6854a438f6ab75b1ffb2fbdc4ccc840aaad50400000000fdffffff02a6800400000000002200203753284c2a9644f2ae9cda8c3fa7433c157755f2684eeec17cd00abe5d7541afa9b4200000000000160014b84ef1e7e398d56fa88a356df3139ff997114f040247304402202b88a9317decf5e219ed44f6947ca28fec4e2ad378496d5b5d6196a7c253518d022017503642ecd2b2b451e68674d1ea2b2006432809cd32d613a59e21221b49b3580121033b3b26324525c22e39fe7942bcac87976804b776c841bbf51f05d22e75a269c700000000","blockhash":"00000000000000000001246cadf2834cf70fe92404e37c14e071f4b7da61993d","confirmations":713,"time":1739204667,"blocktime":1739204667},"id":"getblock.io"}`
	bestblockResponse   = `{"jsonrpc":"2.0","result":"000000000000000000012b7934bb0e26e4ebb12b62290c6407ad252e8d254ffc","id":"getblock.io"}`
)

func TestClient_GetMempool(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rpcReq jsonrpc.RPCRequest[any]
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rpcReq); err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "getrawmempool", rpcReq.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(rawMempoolResponse))
	}))
	defer api.Close()

	client := utxo.NewClient(configs.UtxoRPC{URL: api.URL})

	ctx := t.Context()

	result, err := client.GetMempool(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	assert.Len(t, result, 6)
}

func TestClient_GetBlockByHash(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rpcReq jsonrpc.RPCRequest[any]
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rpcReq); err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "getblock", rpcReq.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(blockResponse))
	}))
	defer api.Close()

	client := utxo.NewClient(configs.UtxoRPC{URL: api.URL})

	ctx := t.Context()

	r, err := client.GetBlockByHash(ctx, "000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506")
	require.NoError(t, err)

	result, ok := r.(*data.Block)
	require.True(t, ok)

	assert.Equal(t, entities.BlockHash("000000000003ba27aa200b1cecaad478d2b00432346c3f1f3986da1afd33e506"), result.Hash)
	assert.NotEmpty(t, result.Tx)
	assert.Equal(t, entities.TxID("fff2525b8931402dd09222c50775608f75787bd2b87e56995a7bdd30f79702c4"), result.Tx[1].TxID)
	assert.Len(t, result.Tx, 4)
}

func TestClient_GetTransactionByTxId(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rpcReq jsonrpc.RPCRequest[any]
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rpcReq); err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "getrawtransaction", rpcReq.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(transactionResponse))
	}))
	defer api.Close()

	client := utxo.NewClient(configs.UtxoRPC{URL: api.URL})

	ctx := t.Context()

	r, err := client.GetTransactionByTxID(ctx, "18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e")
	require.NoError(t, err)
	result, ok := r.(*data.Transaction)
	require.True(t, ok)

	assert.Equal(t, entities.TxID("18dfd1e6734f66cffa0524f9acb7b4c400ed8a5694680ea8ba4b9b24bb57635e"), result.TxID)

	assert.NotEmpty(t, result.Vin)
	assert.Len(t, result.Vin, 1)

	assert.Equal(t, entities.TxID("d5aa0a84cc4cdcfbb2ffb175abf638a45468fbf8eb9a2d1aa380e6e41ac1030b"), result.Vin[0].Txid)

	assert.NotEmpty(t, result.Vout)
	assert.Len(t, result.Vout, 2)

	assert.Equal(t, entities.Address("bc1qxafjsnp2jez09t5um2xrlf6r8s2hw40jdp8wastu6q9tuht4gxhsqwag3z"), result.Vout[0].ScriptPubKey.Address)
	assert.Equal(t, entities.Address("bc1qhp80relrnr2kl2y2x4klxyullxt3zncyj9nd3c"), result.Vout[1].ScriptPubKey.Address)
}

func TestClient_GetLatestBlockHash(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rpcReq jsonrpc.RPCRequest[any]
		decoder := json.NewDecoder(r.Body)

		if err := decoder.Decode(&rpcReq); err != nil {
			assert.NoError(t, err)
		}

		assert.Equal(t, "getbestblockhash", rpcReq.Method)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(bestblockResponse))
	}))
	defer api.Close()

	client := utxo.NewClient(configs.UtxoRPC{URL: api.URL})

	ctx := t.Context()

	result, err := client.GetLatestBlockHash(ctx)
	require.NoError(t, err)

	assert.Equal(t, entities.BlockHash("000000000000000000012b7934bb0e26e4ebb12b62290c6407ad252e8d254ffc"), result)
}
