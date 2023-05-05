pragma solidity ^0.5.0;

import "./U256Cumulative.sol";

contract CumulativeU256Test {
    U256 cumulative ;

    constructor() public {    
        cumulative = new U256(1, 100);  // [1, 100]
        require(cumulative.add(99));
        cumulative.sub(99);
        require(cumulative.get() == 99);

        cumulative.add(1);
        require(cumulative.get() == 100);

        cumulative.sub(100);
        require(cumulative.get() == 100);

        cumulative.sub(99);
        require(cumulative.get() == 1);


        cumulative = new U256(0, 100);  // [1, 100]
        require(cumulative.get() == 0);

        require(cumulative.add(99));
        require(cumulative.get() == 99);
        require(cumulative.sub(99));
        // require(cumulative.get() == 0);

        // cumulative.add(1);
        // require(cumulative.get() == 1);

        // cumulative.sub(100);
        // require(cumulative.get() == 100);

        // cumulative.sub(99);
        // require(cumulative.get() == 1);
    }
}