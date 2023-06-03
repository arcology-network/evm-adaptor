// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Int256.sol";

contract Int256N {
    Int256 array = new Int256();
    constructor  (uint length, int256 initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(int256)  {return array.get(idx);}
    function set(uint256 idx, int256 elem) public { array.set(idx, elem); }
}
