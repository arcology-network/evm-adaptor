// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract Bytes is Base {
    constructor() Base(address(0x84)) {}

    function push(bytes memory elem) public virtual { //9e c6 69 25
        Base.pushBack(elem);
    }    

    function pop() public virtual returns(bytes memory) { // 80 26 32 97
        return Base.popBack();  
    }

    function get(uint256 idx) public virtual returns(bytes memory)  { // 31 fe 88 d0
        return Base.getElem(idx);  
    }

    function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        Base.setElem(idx, abi.encodePacked(elem));
    }
}
