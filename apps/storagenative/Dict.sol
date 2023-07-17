// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../../lib//u256/U256.sol";

contract Dict {
    U256 keys = new U256();
    mapping(uint256 => string) private dict;
    function length() public returns(uint256) { return keys.length();}
    function get(uint256 key) public returns(bool, string)  { return keys.get(key);}
    function set(uint256 key, string val) public { dict[key] = val; }

    function destruct() private {
        for (uint i = 0; i < keys.size(); i ++ ) {
            dict[i] = 0;
            Storage.purge(); // Remove the KV pair completely
        }        
    }
}
