package tests

import (
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	cachedstorage "github.com/arcology-network/common-lib/cachedstorage"
	commonlib "github.com/arcology-network/common-lib/common"
	"github.com/arcology-network/concurrenturl/v2"
	urlcommon "github.com/arcology-network/concurrenturl/v2/common"
	curstorage "github.com/arcology-network/concurrenturl/v2/storage"
	urltype "github.com/arcology-network/concurrenturl/v2/type"
	"github.com/arcology-network/concurrenturl/v2/type/commutative"
	evmcommon "github.com/arcology-network/evm/common"
	arbitrator "github.com/arcology-network/urlarbitrator-engine/go-wrapper"
	adaptor "github.com/arcology-network/vm-adaptor/evm"
)

func TestParallelKittiesPerf(t *testing.T) {
	// persistentDB := urlcommon.NewDataStore()
	// meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	// persistentDB.Save(urlcommon.NewPlatform().Eth10Account(), meta)
	// db := urlcommon.NewTransientDB(persistentDB)
	db := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	db.Inject(urlcommon.NewPlatform().Eth10Account(), meta)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(ceoAddress)
	statedb.AddBalance(ceoAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cooAddress)
	statedb.AddBalance(cooAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cfoAddress)
	statedb.AddBalance(cfoAddress, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy KittyCore.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, ceoAddress, 0, coreCodeV2)
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	coreAddress := receipt.ContractAddress
	t.Log(coreAddress)

	// Deploy SaleClockAuction.
	eu, config = prepare(db, 10000001, transitions, []uint32{1})
	transitions, receipt = deploy(eu, config, ceoAddress, 1, saleAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	saleAddress := receipt.ContractAddress
	t.Log(saleAddress)

	// Deploy SiringClockAuction.
	eu, config = prepare(db, 10000002, transitions, []uint32{2})
	transitions, receipt = deploy(eu, config, ceoAddress, 2, siringAuctionCodeV2, coreAddress.Bytes(), []byte{100})
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	sireAddress := receipt.ContractAddress
	t.Log(sireAddress)

	// Deploy GeneScience.
	eu, config = prepare(db, 10000003, transitions, []uint32{3})
	transitions, receipt = deploy(eu, config, ceoAddress, 3, geneScienceCodeV2, []byte{}, coreAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	geneAddress := receipt.ContractAddress
	t.Log(geneAddress)

	var accesses []urlcommon.UnivalueInterface

	// Call setSaleAuctionAddress.
	eu, config = prepare(db, 10000004, transitions, []uint32{4})
	acc, transitions, receipt := runEx(eu, config, &ceoAddress, &coreAddress, 4, true, "setSaleAuctionAddress(address)", saleAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setSiringAuctionAddress.
	eu, config = prepare(db, 10000005, transitions, []uint32{5})
	acc, transitions, receipt = runEx(eu, config, &ceoAddress, &coreAddress, 5, true, "setSiringAuctionAddress(address)", sireAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setGeneScienceAddress.
	eu, config = prepare(db, 10000006, transitions, []uint32{6})
	acc, transitions, receipt = runEx(eu, config, &ceoAddress, &coreAddress, 6, true, "setGeneScienceAddress(address)", geneAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setCOO.
	eu, config = prepare(db, 10000007, transitions, []uint32{7})
	acc, transitions, receipt = runEx(eu, config, &ceoAddress, &coreAddress, 7, true, "setCOO(address)", cooAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	// Call setCFO.
	eu, config = prepare(db, 10000008, transitions, []uint32{8})
	acc, transitions, receipt = runEx(eu, config, &ceoAddress, &coreAddress, 8, true, "setCFO(address)", cfoAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	accesses = append(accesses, acc...)

	conflicts, _, _ := detectConflict(accesses)
	t.Log("AccessRecords\n" + FormatTransitions(accesses))
	if len(conflicts) != 0 {
		t.Error("unexpected conflictions:", conflicts)
		// return
	}

	// Call unpause.
	eu, config = prepare(db, 10000009, transitions, []uint32{9})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 9, true, "unpause()")
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	PrintMemUsage()

	// Call createPromoKitty.
	url = concurrenturl.NewConcurrentUrl(db)
	url.Import(transitions)
	url.PostImport()
	url.Commit([]uint32{10})
	totalTransitions := []urlcommon.UnivalueInterface{}
	totalAccesses := []urlcommon.UnivalueInterface{}
	txs := []uint32{}
	begin := time.Now()

	config = MainConfig()
	config.Coinbase = &coinbase
	config.BlockNumber = new(big.Int).SetUint64(10000010)
	config.Time = new(big.Int).SetUint64(10000010)
	eu = adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, nil, nil, nil, nil)
	url = concurrenturl.NewConcurrentUrl(nil)
	for i := 0; i < 10000; i++ {
		// url = concurrenturl.NewConcurrentUrl(db)
		url.Init(db)
		api := adaptor.NewAPI(db, url)
		statedb := adaptor.NewStateDB(api, db, url)

		// config := MainConfig()
		// config.Coinbase = &coinbase
		// config.BlockNumber = new(big.Int).SetUint64(10000010)
		// config.Time = new(big.Int).SetUint64(10000010)

		// eu := adaptor.NewEU(config.ChainConfig, *config.VMConfig, config.Chain, statedb, api, db, url)
		eu.SetContext(statedb, api, db, url)
		accesses, transitions, receipt = runEx(eu, config, &cooAddress, &coreAddress, uint64(i), false, "createPromoKitty(uint256,address)", []byte{byte((i + 1) / 65536), byte((i + 1) / 256), byte((i + 1) % 256)}, []byte{byte(i / 65536), byte((i + 1) / 256), byte((i + 1) % 256)})
		// _, transitions, receipt = runEx(eu, config, &cooAddress, &coreAddress, uint64(i), false, "createPromoKitty(uint256,address)", []byte{byte((i + 1) / 65536), byte((i + 1) / 256), byte((i + 1) % 256)}, []byte{byte(i / 65536), byte((i + 1) / 256), byte((i + 1) % 256)})
		if receipt.Status != 1 {
			t.Log(receipt)
			t.Fail()
			return
		}

		bs, _ := commonlib.GobEncode(transitions)
		var ts []urlcommon.UnivalueInterface
		commonlib.GobDecode(bs, &ts)
		totalTransitions = append(totalTransitions, ts...)
		// if i == 0 || i == 256 {
		// 	t.Log("\n" + FormatTransitions(accesses))
		// }
		// totalTransitions = append(totalTransitions, transitions...)
		totalAccesses = append(totalAccesses, accesses...)
		txs = append(txs, uint32(i+1))

		if (i+1)%5000 == 0 {
			// t.Log("\n" + FormatTransitions(transitions))
			// t.Log(receipt)
			t.Log("i = ", i, "----------------------------------------------------------------")
			t.Log("time for exec: ", time.Since(begin))

			// begin = time.Now()
			// for j := 0; j < 100; j++ {
			// 	bs, _ := commonlib.GobEncode(totalTransitions[j*500 : (j+1)*500])
			// 	var ts []urlcommon.UnivalueInterface
			// 	commonlib.GobDecode(bs, ts)
			// }
			// t.Log("time for codec: ", time.Since(begin))

			begin = time.Now()
			conflicts, _, _ := detectConflict(totalAccesses)
			if len(conflicts) != 0 {
				t.Error("unexpected conflicts:", conflicts)
				return
			}
			t.Log("time for detect conflicts: ", time.Since(begin))

			begin = time.Now()
			url = concurrenturl.NewConcurrentUrl(db)
			url.Import(totalTransitions)
			url.PostImport()
			url.Commit(txs)
			t.Log("time for commit: ", time.Since(begin))
			begin = time.Now()
			totalTransitions = []urlcommon.UnivalueInterface{}
			totalAccesses = []urlcommon.UnivalueInterface{}
			txs = []uint32{}
			PrintMemUsage()
		}
	}
}

var engine = arbitrator.Start()

func TestParallelKittiesTransfer(t *testing.T) {
	persistentDB := cachedstorage.NewDataStore()
	meta, _ := commutative.NewMeta(urlcommon.NewPlatform().Eth10Account())
	persistentDB.Inject(urlcommon.NewPlatform().Eth10Account(), meta)
	db := curstorage.NewTransientDB(persistentDB)

	url := concurrenturl.NewConcurrentUrl(db)
	api := adaptor.NewAPI(db, url)
	statedb := adaptor.NewStateDB(api, db, url)
	statedb.Prepare(evmcommon.Hash{}, evmcommon.Hash{}, 0)
	statedb.CreateAccount(coinbase)
	statedb.CreateAccount(ceoAddress)
	statedb.AddBalance(ceoAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cooAddress)
	statedb.AddBalance(cooAddress, new(big.Int).SetUint64(1e18))
	statedb.CreateAccount(cfoAddress)
	statedb.AddBalance(cfoAddress, new(big.Int).SetUint64(1e18))
	_, transitions := url.Export(true)

	// Deploy KittyCore.
	eu, config := prepare(db, 10000000, transitions, []uint32{0})
	transitions, receipt := deploy(eu, config, ceoAddress, 0, coreCode)
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	coreAddress := receipt.ContractAddress
	t.Log(coreAddress)

	// Deploy SaleClockAuction.
	eu, config = prepare(db, 10000001, transitions, []uint32{1})
	transitions, receipt = deploy(eu, config, ceoAddress, 1, saleCode, coreAddress.Bytes(), []byte{100})
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	saleAddress := receipt.ContractAddress
	t.Log(saleAddress)

	// Deploy SiringClockAuction.
	eu, config = prepare(db, 10000002, transitions, []uint32{2})
	transitions, receipt = deploy(eu, config, ceoAddress, 2, sireCode, coreAddress.Bytes(), []byte{100})
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	sireAddress := receipt.ContractAddress
	t.Log(sireAddress)

	// Deploy GeneScience.
	eu, config = prepare(db, 10000003, transitions, []uint32{3})
	transitions, receipt = deploy(eu, config, ceoAddress, 3, geneCode, []byte{}, coreAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	geneAddress := receipt.ContractAddress
	t.Log(geneAddress)

	// Call setSaleAuctionAddress.
	eu, config = prepare(db, 10000004, transitions, []uint32{4})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 4, true, "setSaleAuctionAddress(address)", saleAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call setSiringAuctionAddress.
	eu, config = prepare(db, 10000005, transitions, []uint32{5})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 5, true, "setSiringAuctionAddress(address)", sireAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call setGeneScienceAddress.
	eu, config = prepare(db, 10000006, transitions, []uint32{6})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 6, true, "setGeneScienceAddress(address)", geneAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call setCOO.
	eu, config = prepare(db, 10000007, transitions, []uint32{7})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 7, true, "setCOO(address)", cooAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call setCFO.
	eu, config = prepare(db, 10000008, transitions, []uint32{8})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 8, true, "setCFO(address)", cfoAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call unpause.
	eu, config = prepare(db, 10000009, transitions, []uint32{9})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 9, true, "unpause()")
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)

	// Call createPromoKitty. Assign a kitty to ceo.
	eu, config = prepare(db, 10000010, transitions, []uint32{10})
	transitions, receipt = run(eu, config, &cooAddress, &coreAddress, 0, true, "createPromoKitty(uint256,address)", []byte{1}, ceoAddress.Bytes())
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
	kittyId := receipt.Logs[0].Data[:32]
	t.Log(kittyId)

	// Call transfer. Transfer ceo's kitty to cfo.
	eu, config = prepare(db, 10000011, transitions, []uint32{1})
	transitions, receipt = run(eu, config, &ceoAddress, &coreAddress, 10, true, "transfer(address,uint256)", cfoAddress.Bytes(), kittyId)
	t.Log("\n" + FormatTransitions(transitions))
	t.Log(receipt)
}

func detectConflict(transitions []urlcommon.UnivalueInterface) ([]uint32, []uint32, []bool) {
	start := time.Now()
	length := len(transitions)
	txs := make([]uint32, length)
	paths := make([]string, length)
	reads := make([]uint32, length)
	writes := make([]uint32, length)
	composite := make([]bool, length)
	uniqueTxsDict := make(map[uint32]struct{})
	for i, t := range transitions {
		txs[i] = t.(*urltype.Univalue).GetTx()
		paths[i] = *(t.(*urltype.Univalue).GetPath())
		reads[i] = t.(*urltype.Univalue).Reads()
		writes[i] = t.(*urltype.Univalue).Writes()
		composite[i] = t.(*urltype.Univalue).Composite()
		uniqueTxsDict[txs[i]] = struct{}{}
	}

	uniqueTxs := make([]uint32, 0, len(uniqueTxsDict))
	for tx := range uniqueTxsDict {
		uniqueTxs = append(uniqueTxs, tx)
	}
	fmt.Println(time.Since(start))

	// engine := arbitrator.Start()
	NumPerBatch := 13500
	for i := 0; i < len(txs)/NumPerBatch; i++ {
		start := time.Now()
		arbitrator.Insert(engine, txs[i*NumPerBatch:(i+1)*NumPerBatch], paths[i*NumPerBatch:(i+1)*NumPerBatch], reads[i*NumPerBatch:(i+1)*NumPerBatch], writes[i*NumPerBatch:(i+1)*NumPerBatch], composite[i*NumPerBatch:(i+1)*NumPerBatch])
		fmt.Println(i, time.Since(start))
	}

	// start := time.Now()
	// arbitrator.Insert(engine, txs, paths, reads, writes, composite)
	// fmt.Println(time.Since(start))
	start = time.Now()
	txs, groups, flags := arbitrator.DetectLegacy(engine, uniqueTxs)
	fmt.Println(time.Since(start))
	start = time.Now()
	arbitrator.Clear(engine)
	fmt.Println(time.Since(start))
	return txs, groups, flags
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
