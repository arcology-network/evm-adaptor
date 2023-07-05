// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../noncommutative/bytes/Bytes.sol";

contract ParallelProcess is Base {    
    // Start processing all the parallel jobs in the queue by specifing the number of threads, the number is between [1, 255]
    function run() public returns(bool) {
        (bool success,) = address(0x90).call(abi.encodeWithSignature("run()"));   
        return success;
    } 

    function at(uint256 idx) public returns(bytes memory)  { // 31 fe 88 d0
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("at(uint256)", idx));
        require(success);
        return data;  
    }
}
