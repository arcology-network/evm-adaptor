pragma solidity ^0.5.0;

import "./Invoker.sol";

contract ParallelInvokeTest {
    function callPara() public  {
      

        bytes memory byteArray3 = new bytes(5);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray3[i] = 0x52;
        }

        bytes memory byteArray = new bytes(15);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray[i] = 0x52;
        }

       bytes memory callArg = abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray);
       (bool success0, bytes memory id) = address(this).call(callArg);
       assert(success0);

       addJobTester(abi.encode(address(this), callArg));

       (success0, id) =address(address(0x90)).call(abi.encodeWithSignature("addJob(address,bytes)", address(this), callArg));   
       assert(success0);
    }

    function localTester (address addr, bytes memory args) public returns(bool, bytes memory)  {
        (bool success, bytes memory id) = address(addr).call(args);
        return (success, id);
    }

    function addJobTester (bytes memory encoded) public returns(bool, bytes memory)  {
        (address addr, bytes memory funcCall) = abi.decode(encoded, (address,bytes));
        (bool success0, bytes memory id) = addr.call(funcCall);        
        assert(success0);
        assert(abi.decode(id, (uint256)) == 112);        
        return (true, funcCall);
    }

    function jobExample(address addr, bytes memory id1) pure public returns(uint256){
      return 112;
    }
}