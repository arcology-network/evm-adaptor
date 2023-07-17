// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract U256 is Base {
    constructor() {}

    function push(uint256 elem) public virtual{ //9e c6 69 25
       Base.setKey(uuid(), abi.encode(elem));
    }    

    function pop() public virtual returns(uint256) { // 80 26 32 97
        return abi.decode(Base.popBack(), (uint256));  
    }

    function get(uint256 idx) public virtual returns(uint256)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (uint256));  
    }

    function set(uint256 idx, uint256 elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));
    }
}
