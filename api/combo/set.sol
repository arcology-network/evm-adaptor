// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../noncommutative/base/Base.sol";

contract U256Set is Base { 

   constructor() {}

    function exist(uint256 key) public virtual returns(bool) { //9e c6 69 25
        return get(key) != type(uint256).max;
    }

    function set(uint256 key) public { // 80 26 32 97
        Base.setKey((abi.encodePacked(key)), "0");       
    }

    function get(uint256 key) public virtual returns(uint256){ //9e c6 69 25
        bytes memory data = Base.getKey(abi.encode(key));
        if (data.length > 0) {
           return abi.decode(data, (uint256));  
        }
       return type(uint256).max;
    }    

    function del(uint256 key) public { // 80 26 32 97
        Base.setKey(Base.rand(), abi.encode(key));  
    }

    function length(uint256 idx) public virtual  returns(bool)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (bool));  
    }
}
