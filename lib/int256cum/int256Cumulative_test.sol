// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Int256Cumulative.sol";

contract Int256CumulativeTest {
    Int256Cumulative cumulative ;

    constructor() {    
        cumulative = new Int256Cumulative(0, 1);  // [1, 100]
        // require(cumulative.add(99));
        // cumulative.sub(99);
        // require(cumulative.get() == 99);

        // cumulative.add(1);
        // require(cumulative.get() == 100);

        // cumulative.sub(100);
        // require(cumulative.get() == 100);

        // cumulative.sub(99);
        // require(cumulative.get() == 1);


        // cumulative = new Int256Cumulative(0, 100);  // [1, 100]
        // require(cumulative.get() == 0);

        // require(cumulative.add(99));
        // require(cumulative.get() == 99);
        // require(cumulative.sub(99));
        // require(cumulative.get() == 0);
    }
}