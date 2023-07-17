// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";
import "../runtime/Runtime.sol";

contract Bytes32 is Base {
    constructor() {}

    function push(bytes32 elem) public virtual{ //9e c6 69 25
        Base.setKey(uuid(), abi.encode(elem));
    }    

    function pop() public virtual returns(bytes32) { // 80 26 32 97
        return abi.decode(Base.popBack(), (bytes32));  
    }

    function get(uint256 idx) public virtual returns(bytes32)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (bytes32));  
    }

    function set(uint256 idx, bytes32 elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));
    }
}
