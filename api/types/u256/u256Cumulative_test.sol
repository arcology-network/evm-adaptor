pragma solidity ^0.5.0;

import "./U256Cumulative.sol";

contract CumulativeU256Test {
    U256 cumulative ;

    constructor() public {    
        cumulative = new U256(1, 1, 100);  // [1, 100]
        require(cumulative.add(99));
        require(cumulative.sub(99));
        require(cumulative.get() == 1);

        // require(cumulative.add(uint256(10)));
        // // require(cumulative.get() == 11);
        // require(cumulative.sub(5));
        // // require(cumulative.get() == 6);
        // require(cumulative.sub(0));
        // //  require(cumulative.get() == 6);

        // require(cumulative.add(94));
        // // require(cumulative.get() == 100);
        // // require(!cumulative.add(1));
        // // // require(cumulative.get() == 100);

        // require(cumulative.sub(5));
        // // require(cumulative.get() == 95);
        // require(!cumulative.sub(95));

        // // require(cumulative.get() == 95);
        // require(cumulative.sub(94));
        // require(cumulative.get() == 1);
        // require(cumulative.sub(100));
        // require(cumulative.get() == 0);
        // cumulative.add(uint256(20));

        // cumulative.sub(uint256(2));
        // cumulative.sub(uint256(3));
        // uint256 v = cumulative.get();
        // require(v == 25);
    }
}