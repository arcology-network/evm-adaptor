// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Bytes.sol";

contract BytesN {
    Bytes array = new Bytes();
    constructor (uint256 size, bytes memory initialV) {  
        for (uint i = 0; i < size; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(bytes memory)  {return array.get(idx);}
    function set(uint256 idx, bytes memory elem) public { array.set(idx, elem); }
}
