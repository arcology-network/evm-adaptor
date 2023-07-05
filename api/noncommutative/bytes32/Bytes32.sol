// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract Bytes32 {
    Base base;

    event logMsg(string message);

    constructor() {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(bytes32) { 
        return abi.decode(base.pop(), (bytes32));  
    }

    function push(bytes32 elem) public { 
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(bytes32)  { 
        return abi.decode(base.get(idx), (bytes32));
    }

    function set(uint256 idx, bytes32 elem) public {
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }
}

