// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./U256.sol";

contract U256N {
    U256 array;
    constructor  (uint length, uint256 value) public {  
        array = new U256(); 
        for (uint i = 0; i < length; i ++) {
            array.push(value);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(uint256)  {return array.get(idx);}
    function set(uint256 idx, uint256 elem) public { array.set(idx, elem); }
}
