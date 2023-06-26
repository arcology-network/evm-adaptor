// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./BoolN.sol";

contract BoolNTest {
    BoolN container;

    constructor() {    
        container = new BoolN(4, true);
        require(container.length() == 4); 

        require(container.get(0) == true);
        require(container.get(1) == true);
        require(container.get(2) == true);
        require(container.get(3) == true);

        container.set(0, true);
        container.set(1, false);
        container.set(2, true);
        container.set(3, false);

        require(container.get(0) == true);
        require(container.get(1) == false);
        require(container.get(2) == true);
        require(container.get(3) == false);    
    }
}