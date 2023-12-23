# evm-adaptor

EVM-adaptor is a thin adaptor designed to connecting [parallelized EVM](https://github.com/arcology-network/concurrent-evm) to the rest of the Arcology execution related modules.

## Features

The EVM-adaptor comprises several key components to seamlessly interface with Arcology's concurrency state management system:

### StateDB Implementation:

A new [StateDB](https://github.com/ethereum/go-ethereum/blob/master/core/vm/interface.go) implementation is provided, redirecting all state accesses to Arcology's concurrency state management system.

###  Handlers for State Accesses from the Concurrent APIs
The [Concurrent lib](https://github.com/arcology-network/concurrentlib) provides a variety of concurrent containers and [Solidity API](https://doc.arcology.network/arcology-concurrent-programming-guide/overview) interfaces, assisting developers in creating contracts capable of full parallel processing. The EVM adaptor functions as the module that connects concurrent API calls to Arcology's concurrent state management module through a set of handlers.

   - Byte Array Handler
   - Cumulative Uint256 Handler
   - Cumulative Uint64 Handler
   - Runtime Handler
   - IO Handler
   - Multiprocessor Handler

## Documentation

For detailed information on how to use and integrate the EVM-adaptor into your project, refer to the [documentation.](https://doc.arcology.network)

## License
This project is licensed under the MIT License.