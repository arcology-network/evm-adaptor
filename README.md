<h1> evm-adaptor  <img align="center" height="40" src="./img/evm.svg">  </h1>

EVM-adaptor is a thin adaptor designed to connecting [a parallelized EVM](https://github.com/arcology-network/concurrent-evm) instance to the rest of the Arcology execution related modules.

## Features

The EVM-adaptor comprises several key components to interface with Arcology's concurrency state management system.

### 1. StateDB Implementation:

A new [StateDB](https://github.com/ethereum/go-ethereum/blob/master/core/vm/interface.go) implementation is provided, redirecting all state accesses to Arcology's concurrency state management system.

### 2. Concurrent Container Handler
The [Concurrent lib](https://github.com/arcology-network/concurrentlib) provides a variety of concurrent containers and [Solidity API](https://doc.arcology.network/arcology-concurrent-programming-guide/overview) interfaces, assisting developers in creating contracts capable of full parallel processing. The EVM adaptor functions as the module that connects concurrent API calls to Arcology's concurrent state management module through a set of handlers.

   - Byte Array Handler
   - Cumulative Uint256 Handler
   - Cumulative Uint64 Handler
   - Runtime Handler
   - IO Handler
 
### 3. Multiprocessor Handler

There is a package in the libraray called `Multiprocessor`, where users can start multiple threads, similar to what they can do in other general-purpose languages. The multiprocessor handler takes care of everything related to the `Multiprocessor` contract.

```solidty
   Multiprocessor jobs = Multiprocessor(2);
   
   Jobs.Push(100000, address(this), abi.encodeWithSignature("mint(address,uint256)"), Alice, 100);
   Jobs.Push(100000, address(this), abi.encodeWithSignature("mint(address,uint256)"), Bob, 100);
   
   Jobs.run();
```

## Documentation

For detailed information on how to use and integrate the EVM-adaptor into your project, refer to the [documentation.](https://doc.arcology.network)

## License
This project is licensed under the MIT License.