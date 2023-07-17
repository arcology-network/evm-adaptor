// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract Address is Base {
    constructor() {}

    function push(address elem) public virtual{ //9e c6 69 25
        Base.setKey(uuid(), abi.encode(elem));
    }    

    function pop() public virtual returns(address) { // 80 26 32 97
        return abi.decode(Base.popBack(), (address));  
    }

    function get(uint256 idx) public virtual returns(address)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (address));  
    }

    function set(uint256 idx, address elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));
    }
}
