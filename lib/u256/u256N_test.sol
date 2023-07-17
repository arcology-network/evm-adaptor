// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./U256N.sol";

contract U256NTest {
    U256N container;

    constructor() {    
        container = new U256N(4, 0);
        require(container.length() == 4); 

        require(container.get(0) == uint256(0));
        require(container.get(1) == uint256(0));
        require(container.get(2) == uint256(0));
        require(container.get(3) == uint256(0));

        container.set(0, uint256(11));
        container.set(1, uint256(12));
        container.set(2, uint256(13));
        container.set(3, uint256(14));

        require(container.get(0) == uint256(11));
        require(container.get(1) == uint256(12));
        require(container.get(2) == uint256(13));
        require(container.get(3) == uint256(14));    
    }
}