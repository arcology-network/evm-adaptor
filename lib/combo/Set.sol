// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

// import "../noncommutative/base/Base.sol";
import "./Map.sol";

contract U256Set is U256Map { 
   constructor() {}

    function set(uint256 key) public { // 80 26 32 97
        Base.setKey((abi.encodePacked(key)), abi.encodePacked(uint256(1)));       
    }
}
