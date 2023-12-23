# EVM-adaptor

EVM-adaptor is a thin wrapper designed to connecting [parallelized EVM](https://github.com/arcology-network/concurrent-evm) to the rest of the Arcology execution related modules.

## Features

The EVM-adaptor comprises several key components to seamlessly interface with Arcology's concurrency state management system:

### StateDB Implementation:

   - A new [StateDB](https://github.com/ethereum/go-ethereum/blob/master/core/vm/interface.go) implementation is provided, redirecting all state accesses to Arcology's concurrency state management system.

###  Handlers for State Accesses from the Concurrent APIs

   - Byte Array Handler
   - Cumulative Uint256 Handler
   - Cumulative Uint64 Handler
   - Runtime Handler
   - IO Handler
   - Multiprocessor Handler

## Getting Started

To get started with the EVM-adaptor, follow these steps:

## Documentation

For detailed information on how to use and integrate the EVM-adaptor into your project, refer to the [documentation.](https://doc.arcology.network)

## License
This project is licensed under the MIT License.