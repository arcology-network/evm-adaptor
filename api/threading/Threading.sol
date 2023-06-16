// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Threading {
    address constant private API = address(0x90); 
    bytes private id;

    constructor (uint256 numThreads) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new(uint256)", numThreads));     
        require(success);
        id = data; 
    }

    // Append a new task to the queue, the execution only starts when run() is called
    function add(uint256 gaslimit, address addr, bytes memory args) public returns(bool)  {
        (bool success,) = address(API).call(abi.encodeWithSignature("add(bytes,uint256,address,bytes)", id, gaslimit, addr, args));   
        return (success);
    }

    // Return the size of the queue
    function length() public returns(uint256)  {
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("length(bytes)", id));   
        return abi.decode(data, (uint256));
    }

    // Start processing all the parallel jobs in the queue by specifing the number of threads, the number is between [1, 255]
    function run() public returns(bool) {
        (bool success,) = address(API).call(abi.encodeWithSignature("run(bytes)", id));   
        return success;
    } 
    
    // Get the return data by index.
    function get(uint256 idx) public returns(bool, bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(bytes,uint256)", id, idx));   
        return (success, data);
    }

    // Clear queue and set the length back to zero. Returns true if successful and false otherwise.
    function clear() public returns(bool)  {
        (bool success,) = address(API).call(abi.encodeWithSignature("clear(bytes)", id));   
        return success;
    }

    //Return True if the queue is empty, False otherwise. 
    function  empty() public returns(bool)  {
        return length() == 0;
    }
}
