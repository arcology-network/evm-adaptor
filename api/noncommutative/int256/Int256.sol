// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract Int256 is Base {
    constructor() Base(address(0x84)) {}

    function push(int256 elem) public virtual{ //9e c6 69 25
        Base.setKey(Base.rand(), abi.encode(elem));
    }    

    function pop() public virtual returns(int256) { // 80 26 32 97
        return abi.decode(Base.popBack(), (int256));  
    }

    function get(uint256 idx) public virtual returns(int256)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (int256));  
    }

    function set(uint256 idx, int256 elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));  
    }
}
