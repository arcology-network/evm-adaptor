// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Int256.sol";

contract Int256Test {
    Int256 container = new Int256();
    
    constructor() {     
       require(container.length() == 0); 
    
        container.push((10));
        container.push((-20));
        container.push((30));
        container.push((40));
        require(container.length() == 4); 

        require(container.get(0) == (10));
        require(container.get(1) == (-20));
        require(container.get(2) == (30));
        require(container.get(3) == (40));    

        container.set(0, (-11));
        container.set(1, (12));
        container.set(2, (13));
        container.set(3, (14));

        require(container.get(0) == (-11));
        require(container.get(1) == (12));
        require(container.get(2) == (13));
        require(container.get(3) == (14));

        require(container.pop() == (14));
        require(container.pop() == (13));
        require(container.pop() == (12));
        require(container.pop() == (-11));
        require(container.length() == 0); 
    }
}