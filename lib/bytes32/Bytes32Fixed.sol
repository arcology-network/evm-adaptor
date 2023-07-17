// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Bytes32.sol";

contract Bytes32N {
    Bytes32 array = new Bytes32();
    constructor  (uint length, bytes32 initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(bytes32)  {return array.get(idx);}
    function set(uint256 idx, bytes32 elem) public { array.set(idx, elem); }
}
