// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../noncommutative/base/Base.sol";

contract U256Set is Base {
   mapping (uint256 => uint256) map;
    
   constructor() Base(address(0x84)) {}

   function exist(uint256 key) public virtual returns(bool) { //9e c6 69 25
        return Base.find(abi.encode(key)) < type(uint256).max;
    }

    function insert(uint256 elem) public { // 80 26 32 97
        Base.insert(Base.rand(), abi.encode(elem));  
    }

    function get(uint256 elem) public virtual { //9e c6 69 25
         return abi.decode(Base.getKey(idx), (bool));  
    }    

    function length(uint256 idx) public virtual  returns(bool)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (bool));  
    }
}
