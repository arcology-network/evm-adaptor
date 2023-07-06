// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract String {
    Base base;

    event logMsg(string message);

    constructor() {  base = new Base(address(0x84)); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(string memory) { 
        return abi.decode(base.pop(), (string));  
    }

    function push(string memory elem) public { 
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(string memory)  { 
        return abi.decode(base.get(idx), (string));
    }

    function set(uint256 idx, string memory elem) public {
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }
}


