package tests

import (
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl"
	"github.com/arcology-network/concurrenturl/commutative"
	"github.com/arcology-network/concurrenturl/interfaces"
	curstorage "github.com/arcology-network/concurrenturl/storage"
	evmcommon "github.com/arcology-network/evm/common"
	arbitrator "github.com/arcology-network/urlarbitrator-engine/go-wrapper"
	vmadaptor "github.com/arcology-network/vm-adaptor"
	ccApi "github.com/arcology-network/vm-adaptor/api"
	eucommon "github.com/arcology-network/vm-adaptor/common"
	eth "github.com/arcology-network/vm-adaptor/eth"
	"github.com/arcology-network/vm-adaptor/tests"
)

var (
	coinbase   = eucommon.Coinbase
	ceoAddress = eucommon.Owner
	cooAddress = eucommon.Alice
	cfoAddress = eucommon.Bob
)

func TestParallelKittiesPerf(t *testing.T) {

	db := cachedstorage.NewDataStore()
	meta := commutative.NewPath()
	db.Inject((&concurrenturl.Platform{}).Eth10Account(), meta)

	url := concurrenturl.NewConcurrentUrl(db)
	api := ccApi.NewAPI(url)
	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(ceoAddress)
	statedb.AddBalance(ceoAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cooAddress)
	statedb.AddBalance(cooAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cfoAddress)
	statedb.AddBalance(cfoAddress, new(big.Int).SetUint64(1e18))
	_, transitions := url.ExportAll()

	// Deploy KittyCore.
	eu, config := tests.Prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt, err := tests.Deploy(eu, config, ceoAddress, 0, coreCodeV2)
	if err != nil {
		t.Error("Error: Deploy KittyCore failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	coreAddress := receipt.ContractAddress
	t.Log(coreAddress)

	// Deploy SaleClockAuction.
	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{1})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 1, saleAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	if err != nil {
		t.Error("Error: Deploy SaleClockAuction failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	saleAddress := receipt.ContractAddress
	t.Log(saleAddress)

	// Deploy SiringClockAuction.
	eu, config = tests.Prepare(db, 10000002, transitions, []uint32{2})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 2, siringAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	if err != nil {
		t.Error("Error: Deploy SiringClockAuction failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	sireAddress := receipt.ContractAddress
	t.Log(sireAddress)

	// Deploy GeneScience.
	eu, config = tests.Prepare(db, 10000003, transitions, []uint32{3})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 3, geneScienceCodeV2, []byte{}, coreAddress.Bytes())
	if err != nil {
		t.Error("Error: Deploy GeneScience failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	geneAddress := receipt.ContractAddress
	t.Log(geneAddress)

	var accesses []interfaces.Univalue

	// Call setSaleAuctionAddress.
	eu, config = tests.Prepare(db, 10000004, transitions, []uint32{4})
	acc, transitions, receipt := tests.RunEx(eu, config, &ceoAddress, &coreAddress, 4, true, "setSaleAuctionAddress(address)", saleAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setSiringAuctionAddress.
	eu, config = tests.Prepare(db, 10000005, transitions, []uint32{5})
	acc, transitions, receipt = tests.RunEx(eu, config, &ceoAddress, &coreAddress, 5, true, "setSiringAuctionAddress(address)", sireAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setGeneScienceAddress.
	eu, config = tests.Prepare(db, 10000006, transitions, []uint32{6})
	acc, transitions, receipt = tests.RunEx(eu, config, &ceoAddress, &coreAddress, 6, true, "setGeneScienceAddress(address)", geneAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setCOO.
	eu, config = tests.Prepare(db, 10000007, transitions, []uint32{7})
	acc, transitions, receipt = tests.RunEx(eu, config, &ceoAddress, &coreAddress, 7, true, "setCOO(address)", cooAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setCFO.
	eu, config = tests.Prepare(db, 10000008, transitions, []uint32{8})
	acc, transitions, receipt = tests.RunEx(eu, config, &ceoAddress, &coreAddress, 8, true, "setCFO(address)", cfoAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	conflicts, _, _ := eucommon.DetectConflict(accesses)
	t.Log("AccessRecords\n" + eucommon.FormatTransitions(accesses))
	if len(conflicts) != 0 {
		t.Error("unexpected conflictions:", conflicts)
		// return
	}

	// Call unpause.
	eu, config = tests.Prepare(db, 10000009, transitions, []uint32{9})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 9, true, "unpause()")
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	PrintMemUsage()

	// Call createPromoKitty.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.Sort()
	url.Commit([]uint32{10})
	totalTransitions := []interfaces.Univalue{}
	totalAccesses := []interfaces.Univalue{}
	txs := []uint32{}
	begin := time.Now()

	config = vmadaptor.NewConfig()
	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000010)
	config.Time = new(big.Int).SetUint64(10000010)
	eu = vmadaptor.NewEU(config.ChainConfig, *config.VMConfig, nil, nil)
	url = concurrenturl.NewConcurrentUrl(nil)
	for i := 0; i < 10000; i++ {
		url.Init(db)
		api := ccApi.NewAPI(url)
		statedb := eth.NewImplStateDB(api)
		eu.SetRuntimeContext(statedb, api)
		accesses, transitions, receipt = tests.RunEx(eu, config, &cooAddress, &coreAddress, uint64(i), false, "createPromoKitty(uint256,address)", []byte{byte((i + 1) / 65536), byte((i + 1) / 256), byte((i + 1) % 256)}, []byte{byte(i / 65536), byte((i + 1) / 256), byte((i + 1) % 256)})
		if receipt.Status != 1 {
			t.Log(receipt)
			t.Fail()
			return
		}

		bs, _ := commonlib.GobEncode(transitions)
		var ts []interfaces.Univalue
		commonlib.GobDecode(bs, &ts)
		totalTransitions = append(totalTransitions, ts...)

		totalAccesses = append(totalAccesses, accesses...)
		txs = append(txs, uint32(i+1))

		if (i+1)%5000 == 0 {
			t.Log("i = ", i, "----------------------------------------------------------------")
			t.Log("time for exec: ", time.Since(begin))

			begin = time.Now()
			conflicts, _, _ := eucommon.DetectConflict(totalAccesses)
			if len(conflicts) != 0 {
				t.Error("unexpected conflicts:", conflicts)
				return
			}
			t.Log("time for detect conflicts: ", time.Since(begin))

			begin = time.Now()
			url = concurrenturl.NewConcurrentUrl(db)
			url.Import(totalTransitions)
			url.Sort()
			url.Commit(txs)
			t.Log("time for commit: ", time.Since(begin))
			begin = time.Now()
			totalTransitions = []interfaces.Univalue{}
			totalAccesses = []interfaces.Univalue{}
			txs = []uint32{}
			PrintMemUsage()
		}
	}
}

var engine = arbitrator.Start()

func TestParallelKittiesTransfer(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta := commutative.NewPath()
	persistentDB.Inject((&concurrenturl.Platform{}).Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)

	api := ccApi.NewAPI(url)
	statedb := eth.NewImplStateDB(api)
	statedb.PrepareFormer(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(ceoAddress)
	statedb.AddBalance(ceoAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cooAddress)
	statedb.AddBalance(cooAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cfoAddress)
	statedb.AddBalance(cfoAddress, new(big.Int).SetUint64(1e18))
	_, transitions := url.ExportAll()

	// Deploy KittyCore.
	eu, config := tests.Prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt, err := tests.Deploy(eu, config, ceoAddress, 0, coreCodeV2)
	if err != nil {
		t.Error("Error: Deploy KittyCore failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	coreAddress := receipt.ContractAddress
	t.Log(coreAddress)

	// Deploy SaleClockAuction.
	eu, config = tests.Prepare(db, 10000001, transitions, []uint32{1})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 1, saleAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	if err != nil {
		t.Error("Error: Deploy SaleClockAuction failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	saleAddress := receipt.ContractAddress
	t.Log(saleAddress)

	// Deploy SiringClockAuction.
	eu, config = tests.Prepare(db, 10000002, transitions, []uint32{2})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 2, siringAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	if err != nil {
		t.Error("Error: Deploy SiringClockAuction failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	sireAddress := receipt.ContractAddress
	t.Log(sireAddress)

	// Deploy GeneScience.
	eu, config = tests.Prepare(db, 10000003, transitions, []uint32{3})
	transitions, receipt, err = tests.Deploy(eu, config, ceoAddress, 3, geneScienceCodeV2, []byte{}, coreAddress.Bytes())
	if err != nil {
		t.Error("Error: Deploy GeneScience failed:", err)
		return
	}
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	geneAddress := receipt.ContractAddress
	t.Log(geneAddress)

	// Call setSaleAuctionAddress.
	eu, config = tests.Prepare(db, 10000004, transitions, []uint32{4})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 4, true, "setSaleAuctionAddress(address)", saleAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call setSiringAuctionAddress.
	eu, config = tests.Prepare(db, 10000005, transitions, []uint32{5})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 5, true, "setSiringAuctionAddress(address)", sireAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call setGeneScienceAddress.
	eu, config = tests.Prepare(db, 10000006, transitions, []uint32{6})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 6, true, "setGeneScienceAddress(address)", geneAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call setCOO.
	eu, config = tests.Prepare(db, 10000007, transitions, []uint32{7})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 7, true, "setCOO(address)", cooAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call setCFO.
	eu, config = tests.Prepare(db, 10000008, transitions, []uint32{8})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 8, true, "setCFO(address)", cfoAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call unpause.
	eu, config = tests.Prepare(db, 10000009, transitions, []uint32{9})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 9, true, "unpause()")
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)

	// Call createPromoKitty. Assign a kitty to ceo.
	eu, config = tests.Prepare(db, 10000010, transitions, []uint32{10})
	transitions, receipt = tests.Run(eu, config, &cooAddress, &coreAddress, 0, true, "createPromoKitty(uint256,address)", []byte{1}, ceoAddress.Bytes())
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
	kittyId := receipt.Logs[0].Data[:32]
	t.Log(kittyId)

	// Call transfer. Transfer ceo's kitty to cfo.
	eu, config = tests.Prepare(db, 10000011, transitions, []uint32{1})
	transitions, receipt = tests.Run(eu, config, &ceoAddress, &coreAddress, 10, true, "transfer(address,uint256)", cfoAddress.Bytes(), kittyId)
	//t.Log("\n" + eucommon.FormatTransitions(transitions))
	t.Log(receipt)
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v", m.NumGC)
	// fmt.Printf("\tBySize = %v\n", m.BySize)
	fmt.Printf("\tBySize = [")
	for _, entry := range m.BySize {
		if entry.Mallocs-entry.Frees < 50000 {
			continue
		}
		fmt.Printf("{%v, %v} ", entry.Size, entry.Mallocs-entry.Frees)
	}
	fmt.Printf("]\n")
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
