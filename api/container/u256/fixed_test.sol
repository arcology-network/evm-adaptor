pragma solidity ^0.5.0;

import "./Fixed.sol";

contract U256FixedTest {
    U256Fixed container;
    U256[] array;

    constructor() public {    
        container = new U256Fixed(4, 0);
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