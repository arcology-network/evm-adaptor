// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Int256N.sol";

contract Int64NTest {
    Int256N container;

    constructor() {    
        int256 num1 = 0;
        int256 num2 = 1;
        int256 num3 = 2;
        int256 num4 = 3;

        container = new Int256N(4, num1);
        require(container.length() == 4); 
        
        require((container.get(0)) == (num1));
        require((container.get(1)) == (num1));
        require((container.get(2)) == (num1));
        require((container.get(3)) == (num1));

        container.set(0, num4);
        container.set(1, num3);
        container.set(2, num2);
        container.set(3, num1);

        require((container.get(0)) == (num4));
        require((container.get(1)) == (num3));
        require((container.get(2)) == (num2));
        require((container.get(3)) == (num1));    
    }
}