pragma solidity ^0.5.0;

import "./Parallel.sol";

contract ParallelInvokeTest {
    constructor () public {      
        bytes memory myFixedByteArray = "abcd";

        Parallel parallel  = new Parallel();
        parallel.newJob(myFixedByteArray);   
    }
}