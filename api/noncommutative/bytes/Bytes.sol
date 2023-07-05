// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract Bytes {
    Base base;

    event logMsg(string message);

    constructor() {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(bytes memory) { // 80 26 32 97
        return abi.decode(base.pop(), (bytes));  
    }

    function push(bytes memory elem) public { //9e c6 69 25
        base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(bytes memory)  { // 31 fe 88 d0
        return abi.decode(base.get(idx), (bytes));  
    }

    function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));     
    }
}
