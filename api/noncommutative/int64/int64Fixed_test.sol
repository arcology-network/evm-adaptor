pragma solidity ^0.5.0;

import "./Int64Fixed.sol";

contract Int64FixedTest {
    Int64Fixed container;

    constructor() public {    
        int64 num1 = 0;
        int64 num2 = 1;
        int64 num3 = 2;
        int64 num4 = 3;

        container = new Int64Fixed(4, num1);
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