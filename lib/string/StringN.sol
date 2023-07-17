// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./String.sol";

contract StringN {
    String array = new String();
    constructor  (uint length, string memory initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(string memory)  {return array.get(idx);}
    function set(uint256 idx, string memory elem) public { array.set(idx, elem); }
}
