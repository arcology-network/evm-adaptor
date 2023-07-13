// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract Bool {
    Base base;

    constructor() {  base = new Base(address(0x84), false); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(bool) { 
        return abi.decode(base.pop(), (bool));  
    }

    function push(bool elem) public { 
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(bool)  { 
        return abi.decode(base.get(idx), (bool));
    }

    function set(uint256 idx, bool elem) public {
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }

    function clear() public { // 7a fa 6 2 38
        base.clear();
    }
}
