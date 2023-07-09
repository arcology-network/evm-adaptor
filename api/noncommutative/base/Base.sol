// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Base {
    address public immutable API;// = address(0x84);    
    event logMsg(string message);

    constructor (address addr) {
        API = addr;
        (bool success,) = address(API).call(abi.encodeWithSignature("new()"));       
        require(success);
    }

    function length() public returns(uint256) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));
        require(success);
        return abi.decode(data, (uint256));
    }

    // The initial length of the container at the current block height
    function peek() public returns(bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("peek()"));
        require(success);
        return data;  
    } 

    function pop() public returns(bytes memory) { // 80 26 32 97
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("pop()"));
        require(success);
        return abi.decode(data, (bytes)); 
    }

    function push(bytes memory encoded) public { //9e c6 69 25
        (bool success,) = address(API).call(encoded);
        require(success);
    }  

    function get(uint256 idx) public returns(bytes memory)  { // 31 fe 88 d0
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(uint256)", idx));
        require(success);
        return abi.decode(data, (bytes));  
    }

    function set(bytes memory encoded) public { // 7a fa 62 38
        (bool success,) = address(API).call(encoded);
        require(success);
    }

    //Return True if the queue is empty, False otherwise. 
    function  empty() public returns(bool)  {
        return length() == 0;
    }

    // Clear the data
    function clear() public {
        address(API).call(abi.encodeWithSignature("clear()"));       
    }

    function run(uint256 numThreads) public {
        address(API).call(abi.encodeWithSignature("run()", numThreads));       
    }
    
    // function log(bytes memory elem) public { // 7a fa 62 38
    //     address(API).call(abi.encodeWithSignature("log()", elem));     
    // }
}
