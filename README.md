<h1> evm-adaptor  <img align="center" height="50" src="./img/evm.svg">  </h1>

EVM-adaptor is a thin adaptor designed to connecting [a parallelized EVM](https://github.com/arcology-network/concurrent-evm) instance to the rest of the Arcology execution related modules.

<h2> Components <img align="center" height="32" src="./img/components.svg">  </h2>

The EVM-adaptor comprises several key components to interface with Arcology's concurrency state management system.

### StateDB Implementation:

A new [StateDB](https://github.com/ethereum/go-ethereum/blob/master/core/vm/interface.go) implementation is provided, redirecting all state accesses to Arcology's concurrency state management system.

### Concurrent Container Handler
The [Concurrent lib](https://github.com/arcology-network/concurrentlib) provides a variety of concurrent containers and tools in the [Solidity API](https://doc.arcology.network/arcology-concurrent-programming-guide/overview) interfaces, assisting developers in creating contracts capable of full parallel processing. The EVM adaptor functions as the module that **connects concurrent API calls to Arcology's concurrent state management** module through a set of handlers.

   - Byte Array Handler
   - Cumulative Uint256 Handler
   - Cumulative Uint64 Handler
   - Runtime Handler
   - IO Handler
   - **Multiprocessor Handler**
 
<h2> Documentation <img align="center" height="32" src="./img/doc.svg">  </h2>

For detailed information on how to use and integrate the EVM-adaptor into your project, refer to the [documentation.](https://doc.arcology.network)


<h2> License  <img align="center" height="32" src="./img/copyright.svg">  </h2>
## License
This project is licensed under the MIT License.