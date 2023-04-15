package tests

import (
	"bytes"
	"math/big"
	"testing"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	evmtypes "github.com/arcology-network/evm/core/types"
	"github.com/arcology-network/evm/crypto"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
)

func TestUniswapFunctionTest(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(owner)
	statedb.AddBalance(owner, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(user1)
	statedb.AddBalance(user1, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(user2)
	statedb.AddBalance(user2, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy UniswapV2Factory.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, owner, 0, uniswapFactoryCode, owner.Bytes())
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	factoryAddress := receipt.ContractAddress
	t.Log(factoryAddress)

	// Deploy two tokens.
	eu, config = prepare(db, 10000001, transitions, []uint32{1})
	transitions, receipt = deploy(eu, config, owner, 1, erc20TokenCode, []byte("TKA"))
	// t.Log("\n", FormatTransitions(transitions))
	t.Log(receipt)
	token1Address := receipt.ContractAddress
	t.Log(token1Address)

	eu, config = prepare(db, 10000002, transitions, []uint32{2})
	transitions, receipt = deploy(eu, config, owner, 2, erc20TokenCode, []byte("TKB"))
	// t.Log("\n", FormatTransitions(transitions))
	t.Log(receipt)
	token2Address := receipt.ContractAddress
	t.Log(token2Address)

	// Call createPair.
	eu, config = prepare(db, 10000003, transitions, []uint32{3})
	transitions, receipt = run(eu, config, &owner, &factoryAddress, 3, true, "createPair(address,address)", token1Address.Bytes(), token2Address.Bytes())
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	pairAddress := evmcommon.BytesToAddress(receipt.Logs[0].Data[12:32])
	t.Log(pairAddress)

	// Deploy UniswapV2Router02.
	eu, config = prepare(db, 10000004, transitions, []uint32{4})
	transitions, receipt = deploy(eu, config, owner, 4, uniswapRouterCode, factoryAddress.Bytes(), []byte{})
	// t.Log("\n", FormatTransitions(transitions))
	t.Log(receipt)
	routerAddress := receipt.ContractAddress
	t.Log(routerAddress)

	// Mint on TKA and TKB for user1.
	eu, config = prepare(db, 10000005, transitions, []uint32{5})
	transitions, receipt = run(eu, config, &owner, &token1Address, 5, true, "mint(address,uint256)", user1.Bytes(), []byte{1, 0, 0, 0, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	eu, config = prepare(db, 10000006, transitions, []uint32{6})
	transitions, receipt = run(eu, config, &owner, &token2Address, 6, true, "mint(address,uint256)", user1.Bytes(), []byte{1, 0, 0, 0, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// User1 approve UniswapV2Router02 to call transferFrom.
	// This is the preparation for calling addLiquidity.
	eu, config = prepare(db, 10000007, transitions, []uint32{7})
	transitions, receipt = run(eu, config, &user1, &token1Address, 0, true, "approve(address,uint256)", routerAddress.Bytes(), []byte{1, 0, 0, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	eu, config = prepare(db, 10000008, transitions, []uint32{1})
	transitions, receipt = run(eu, config, &user1, &token2Address, 1, true, "approve(address,uint256)", routerAddress.Bytes(), []byte{1, 0, 0, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// User1 call addLiquidity.
	eu, config = prepare(db, 10000009, transitions, []uint32{2})
	transitions, receipt = run(
		eu, config, &user1, &routerAddress, 2, true,
		"addLiquidity(address,address,uint256,uint256,uint256,uint256,address,uint256)",
		token1Address.Bytes(),
		token2Address.Bytes(),
		[]byte{1, 0, 0, 0},
		[]byte{1, 0, 0, 0},
		[]byte{1, 0, 0},
		[]byte{1, 0, 0},
		user1.Bytes(),
		[]byte{1, 0, 0, 0, 0},
	)
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	if !bytes.Equal(receipt.Logs[0].Topics[1][12:], user1.Bytes()) || // token1Address.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[0].Topics[2][12:], pairAddress.Bytes()) || // token1Address.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[0].Data, evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) || // token1Address.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[1].Topics[1][12:], user1.Bytes()) || // token2Address.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[1].Topics[2][12:], pairAddress.Bytes()) || // token2Address.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[1].Data, evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) || // token2Address.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[2].Topics[1][:], evmcommon.Hash{}.Bytes()) || // pairAddress.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[2].Topics[2][:], evmcommon.Hash{}.Bytes()) || // pairAddress.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[2].Data, evmcommon.BytesToHash([]byte{3, 232}).Bytes()) || // pairAddress.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[3].Topics[1][:], evmcommon.Hash{}.Bytes()) || // pairAddress.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[3].Topics[2][12:], user1.Bytes()) || // pairAddress.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[3].Data, evmcommon.BytesToHash([]byte{255, 252, 24}).Bytes()) || // pairAddress.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[4].Data[:32], evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) || // pairAddress.Sync(*reserve0*, reserve1)
		!bytes.Equal(receipt.Logs[4].Data[32:], evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) || // pairAddress.Sync(reserve0, *reserve1*)
		!bytes.Equal(receipt.Logs[5].Topics[1][12:], routerAddress.Bytes()) || // pairAddress.Mint(*sender*, amount0, amount1)
		!bytes.Equal(receipt.Logs[5].Data[:32], evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) || // pairAddress.Mint(sender, *amount0*, amount1)
		!bytes.Equal(receipt.Logs[5].Data[32:], evmcommon.BytesToHash([]byte{1, 0, 0, 0}).Bytes()) { // pairAddress.Mint(sender, amount0, *amount1*)
		t.Fail()
	}

	// Mint on TKA for user2.
	eu, config = prepare(db, 10000010, transitions, []uint32{3})
	transitions, receipt = run(eu, config, &owner, &token1Address, 7, true, "mint(address,uint256)", user2.Bytes(), []byte{1, 0, 0, 0, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// User2 approve UniswapV2Router02 to call transferFrom.
	// This is the preparation for calling swap.
	eu, config = prepare(db, 10000011, transitions, []uint32{8})
	transitions, receipt = run(eu, config, &user2, &token1Address, 0, true, "approve(address,uint256)", routerAddress.Bytes(), []byte{1, 0})
	// t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// User2 call swapExactTokensForTokens.
	eu, config = prepare(db, 10000012, transitions, []uint32{1})
	transitions, receipt = run(
		eu, config, &user2, &routerAddress, 1, true,
		"swapExactTokensForTokens(uint256,uint256,address[],address,uint256)",
		[]byte{1, 0},
		[]byte{1},
		[]byte{0xa0},
		user2.Bytes(),
		[]byte{1, 0, 0, 0, 0},
		[]byte{2},
		token1Address.Bytes(),
		token2Address.Bytes(),
	)
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	if !bytes.Equal(receipt.Logs[0].Topics[1][12:], user2.Bytes()) || // token1Address.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[0].Topics[2][12:], pairAddress.Bytes()) || // token1Address.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[0].Data, evmcommon.BytesToHash([]byte{1, 0}).Bytes()) || // token1Address.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[1].Topics[1][12:], pairAddress.Bytes()) || // token2Address.Transfer(*from*, to, value)
		!bytes.Equal(receipt.Logs[1].Topics[2][12:], user2.Bytes()) || // token2Address.Transfer(from, *to*, value)
		!bytes.Equal(receipt.Logs[1].Data, evmcommon.BytesToHash([]byte{255}).Bytes()) || // token2Address.Transfer(from, to, *value*)
		!bytes.Equal(receipt.Logs[2].Data[:32], evmcommon.BytesToHash([]byte{1, 0, 1, 0}).Bytes()) || // pairAddress.Sync(*reserve0*, reserve1)
		!bytes.Equal(receipt.Logs[2].Data[32:], evmcommon.BytesToHash([]byte{255, 255, 1}).Bytes()) || // pairAddress.Sync(reserve0, *reserve1*)
		!bytes.Equal(receipt.Logs[3].Topics[1][12:], routerAddress.Bytes()) || // pairAddress.Swap(*from*, amount0In, amount1In, amount0Out, amount1Out, to)
		!bytes.Equal(receipt.Logs[3].Topics[2][12:], user2.Bytes()) || // pairAddress.Swap(*from*, amount0In, amount1In, amount0Out, amount1Out, *to*)
		!bytes.Equal(receipt.Logs[3].Data[:32], evmcommon.BytesToHash([]byte{1, 0}).Bytes()) || // pairAddress.Swap(from, *amount0In*, amount1In, amount0Out, amount1Out, to)
		!bytes.Equal(receipt.Logs[3].Data[32:64], evmcommon.Hash{}.Bytes()) || // pairAddress.Swap(from, amount0In, *amount1In*, amount0Out, amount1Out, to)
		!bytes.Equal(receipt.Logs[3].Data[64:96], evmcommon.Hash{}.Bytes()) || // pairAddress.Swap(from, amount0In, amount1In, *amount0Out*, amount1Out, to)
		!bytes.Equal(receipt.Logs[3].Data[96:], evmcommon.BytesToHash([]byte{255}).Bytes()) { // pairAddress.Swap(from, amount0In, amount1In, amount0Out, *amount1Out*, to)
		t.Fail()
	}
}

func prepare(db urlcommon.DatastoreInterface, height uint64, transitions []urlcommon.UnivalueInterface, txs []uint32) (*adaptor.EU, *adaptor.Config) {
	url := concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit(txs)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)

	config := MainTestConfig()
	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(height)
	config.Time = new(big.Int).SetUint64(height)

	return adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url), config
}

func deploy(eu *adaptor.EU, config *adaptor.Config, owner evmcommon.Address, nonce uint64, code string, args ...[]byte) ([]urlcommon.UnivalueInterface, *evmtypes.Receipt) {
	data := evmcommon.Hex2Bytes(code)
	for _, arg := range args {
		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
	}
	msg := evmtypes.NewMessage(owner, nil, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, false)
	_, transitions, receipt := eu.Run(evmcommon.BytesToHash([]byte{byte(nonce + 1), byte(nonce + 1), byte(nonce + 1)}), int(nonce+1), &msg, adaptor.NewEVMBlockContext(config), adaptor.NewEVMTxContext(msg))
	return transitions, receipt
}

func run(eu *adaptor.EU, config *adaptor.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, args ...[]byte) ([]urlcommon.UnivalueInterface, *evmtypes.Receipt) {
	data := crypto.Keccak256([]byte(function))[:4]
	for _, arg := range args {
		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
	}
	msg := evmtypes.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
	_, transitions, receipt := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), int(nonce+1), &msg, adaptor.NewEVMBlockContext(config), adaptor.NewEVMTxContext(msg))
	return transitions, receipt
}

func runEx(eu *adaptor.EU, config *adaptor.Config, from, to *evmcommon.Address, nonce uint64, checkNonce bool, function string, args ...[]byte) ([]urlcommon.UnivalueInterface, []urlcommon.UnivalueInterface, *evmtypes.Receipt) {
	data := crypto.Keccak256([]byte(function))[:4]
	for _, arg := range args {
		data = append(data, evmcommon.BytesToHash(arg).Bytes()...)
	}
	msg := evmtypes.NewMessage(*from, to, nonce, new(big.Int).SetUint64(0), 1e15, new(big.Int).SetUint64(1), data, nil, checkNonce)
	accesses, transitions, receipt := eu.Run(evmcommon.BytesToHash([]byte{byte((nonce + 1) / 65536), byte((nonce + 1) / 256), byte((nonce + 1) % 256)}), int(nonce+1), &msg, adaptor.NewEVMBlockContext(config), adaptor.NewEVMTxContext(msg))
	return accesses, transitions, receipt
}
