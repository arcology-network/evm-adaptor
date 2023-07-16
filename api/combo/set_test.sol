// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./set.sol";

contract SetTest {
    U256Set set = new U256Set();
    constructor() {     
        require(set.length() == 0); 
        set.insert(10);
        set.insert(11);

        require(!set.exist(0));
        require(set.exist(10));
    }
}