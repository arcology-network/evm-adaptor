// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract Int256 {
    Base base;

    constructor() {  base = new Base(address(0x84)); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(int256) { // 80 26 32 97
        return abi.decode(base.pop(), (int256));  
    }

    function push(int256 elem) public{ //9e c6 69 25
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(int256)  { // 31 fe 88 d0
        return abi.decode(base.get(idx), (int256));
    }

    function set(uint256 idx, int256 elem) public { // 7a fa 62 38
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }

    function clear() public { // 7a fa 62 38
        base.clear();
    }
}
