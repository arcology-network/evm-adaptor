pragma solidity ^0.5.0;

import "./Parallel.sol";

contract ParallelInvokeTest {
    function callPara() public returns(int) {
        Parallel parallel  = new Parallel();
        (bool success, uint32 jobID) = parallel.newJob(address(this), abi.encodeWithSignature("jobExample(int) returns(int)", 9999));  
        assert(success); 
        return 65536;
    }

    function jobExample(int counter) pure public returns(int)  {
        return counter ++;
    }
}