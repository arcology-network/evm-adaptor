module github.com/arcology-network/vm-adaptor

go 1.20

replace github.com/arcology-network/evm => ../evm/

replace github.com/arcology-network/common-lib => ../common-lib/

replace github.com/arcology-network/concurrenturl => ../concurrenturl/

replace github.com/arcology-network/urlarbitrator-engine => ../urlarbitrator-engine/

// github.com/arcology-network/common-lib v1.7.1
// github.com/arcology-network/concurrenturl v1.7.1
// require github.com/arcology-network/urlarbitrator-engine v1.7.1

// github.com/arcology-network/evm v0.0.0-20221110011616-1cdc0ab27c8e

require (
	github.com/arcology-network/common-lib v0.0.0-00010101000000-000000000000
	github.com/arcology-network/evm v0.0.0-00010101000000-000000000000
	github.com/arcology-network/urlarbitrator-engine v0.0.0-00010101000000-000000000000
	github.com/holiman/uint256 v1.2.2
)

// require github.com/arcology-network/concurrenturl v1.7.1 // indirect

require (
	github.com/AndreasBriese/bbloom v0.0.0-20190825152654-46b345b51c96 // indirect
	// github.com/HPISTechnologies/common-lib v1.7.0 // indirect
	// github.com/HPISTechnologies/concurrenturl v1.6.2-0.20221103064740-987edfad3b72 // indirect
	// github.com/HPISTechnologies/evm v1.10.4-0.20220902063532-a4690b596277 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/VictoriaMetrics/fastcache v1.6.0 // indirect
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgraph-io/badger v1.6.2 // indirect
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elliotchance/orderedmap v1.5.0 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	golang.org/x/crypto v0.2.0
	golang.org/x/exp v0.0.0-20230510235704-dd950f8aeaea // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
)

// github.com/arcology-network/common-lib v1.7.1
require github.com/arcology-network/concurrenturl v0.0.0-00010101000000-000000000000

require (
	git.sr.ht/~sbinet/gg v0.4.1 // indirect
	github.com/DataDog/zstd v1.5.2 // indirect
<<<<<<< HEAD
=======
	github.com/ajstarks/svgo v0.0.0-20211024235047-1546f124cd8b // indirect
	github.com/arcology-network/3rd-party v1.7.1 // indirect
>>>>>>> f50d5beec9d156a6b0b2497766ddf11d5e3df8ac
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/cockroachdb/errors v1.9.1 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/pebble v0.0.0-20230209160836-829675f94811 // indirect
	github.com/cockroachdb/redact v1.1.3 // indirect
	github.com/deckarep/golang-set/v2 v2.1.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/getsentry/sentry-go v0.18.0 // indirect
<<<<<<< HEAD
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
=======
	github.com/go-fonts/liberation v0.3.1 // indirect
	github.com/go-latex/latex v0.0.0-20230307184459-12ec69307ad9 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-pdf/fpdf v0.8.0 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
>>>>>>> f50d5beec9d156a6b0b2497766ddf11d5e3df8ac
	github.com/klauspost/compress v1.15.15 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.39.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	golang.org/x/image v0.7.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	gonum.org/v1/plot v0.13.0 // indirect
)
