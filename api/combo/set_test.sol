// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Set.sol";

contract SetTest {
    U256Set set = new U256Set();
    constructor() {     
        require(set.length() == 0); 
        set.set(10);
        set.set(11);
        require(set.length() == 2); 

        require(!set.exist(0));
        require(set.exist(10)); 
        require(set.exist(11)); 

        set.del(10);
        require(set.length() == 1); 

        require(set.exist(11));
        set.del(11);
        require(set.length() == 0); 

        require(!set.exist(10));
        require(!set.exist(11));
    }
}