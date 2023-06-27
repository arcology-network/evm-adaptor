// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract U256 {
    Base base;

    constructor() {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(uint256) { // 80 26 32 97
        return abi.decode(abi.decode(base.pop(), (bytes)), (uint256));  
    }

    function push(uint256 elem) public { //9e c6 69 25
       base.push(abi.encodeWithSignature("push(bytes,bytes)",  base.id(), abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(uint256)  { // 31 fe 88 d0
        return abi.decode(abi.decode(base.get(idx), (bytes)), (uint256));
    }

    function set(uint256 idx, uint256 elem) public { // 7a fa 62 38
        base.set(abi.encodeWithSignature("set(bytes,uint256,bytes)", base.id(), idx, abi.encode(elem)));        
    }
}
