// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../base/Base.sol";

contract U256Map is Base { 
   constructor() {}

    function exist(uint256 key) public virtual returns(bool) { //9e c6 69 25
        (bool success,) = get(key);
        return success;
    }

    function set(uint256 key, uint256 value) public { // 80 26 32 97
        Base.setKey((abi.encodePacked(key)), abi.encodePacked(value));       
    }

    function get(uint256 key) public virtual returns(bool, uint256){ //9e c6 69 25
        bytes memory data = Base.getKey(abi.encodePacked(key));
        if (data.length > 0) {
           return (true, abi.decode(data, (uint256)));  
        }
       return (false, type(uint256).max);
    }    

    function del(uint256 key) public { // 80 26 32 97
        Base.delKey((abi.encodePacked(key)));  
    }
}
