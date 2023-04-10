pragma solidity ^0.5.0;

import "./U256.sol";

contract CumulativeU256Test {
    U256 cumulative ;

    constructor() public {    
        cumulative = new U256(0, 1, 100);
        // cumulative.add(uint256(10));
        // cumulative.add(uint256(20));

        // cumulative.sub(uint256(2));
        // cumulative.sub(uint256(3));
        // uint256 v = cumulative.get();
        // require(v == 25);
    }
}