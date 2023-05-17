pragma solidity ^0.5.0;

import "./Int256Cumulative.sol";

contract Int256CumulativeTest {
    Int256Cumulative cumulative ;

    constructor() public {    
        cumulative = new Int256Cumulative(1, 100);  // [1, 100]
        require(cumulative.add(99));
        cumulative.sub(99);
        require(cumulative.get() == 99);

        cumulative.add(1);
        require(cumulative.get() == 100);

        cumulative.sub(100);
        require(cumulative.get() == 100);

        cumulative.sub(99);
        require(cumulative.get() == 1);


        cumulative = new Int256Cumulative(0, 100);  // [1, 100]
        require(cumulative.get() == 0);

        require(cumulative.add(99));
        require(cumulative.get() == 99);
        require(cumulative.sub(99));
        require(cumulative.get() == 0);
    }
}