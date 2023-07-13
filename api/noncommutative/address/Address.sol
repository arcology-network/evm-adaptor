// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract Address {
    Base base;

    event logMsg(string message);

    constructor() {  base = new Base(address(0x84), false); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(address) { 
        return abi.decode(base.pop(), (address));  
    }

    function push(address elem) public { 
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(address)  { 
        return abi.decode(base.get(idx), (address));
    }

    function set(uint256 idx, address elem) public {
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }

    function clear() public { // 7a fa 62 38
        base.clear();
    }
}

