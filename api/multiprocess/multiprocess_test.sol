pragma solidity ^0.5.0;

import "./Multiprocess.sol";

contract MultiprocessTest {
    function call() public  { 
        Multiprocess mp = new Multiprocess();

        bytes memory byteArray = new bytes(15);
        for (uint  i = 0; i < byteArray.length; i ++) {
            byteArray[i] = 0x52;
        }

       mp.addJob(address(this), abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray));
       mp.addJob(address(this), abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray));
       mp.addJob(address(this), abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray));
       mp.addJob(address(this), abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray));
       mp.addJob(address(this), abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray));
       assert(mp.length() == 5);

       (bool success, bytes memory id) = address(address(0x90)).call(abi.encodeWithSignature("run()", address(this)));   
       assert(success);
    }

    function callBasic() public  {      
        bytes memory byteArray3 = new bytes(5);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray3[i] = 0x42;
        }

        bytes memory byteArray = new bytes(15);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray[i] = 0x52;
        }

       bytes memory callArg = abi.encodeWithSignature("jobExample(address,bytes)", address(this),byteArray);
       (bool success0, bytes memory id) = address(this).call(callArg);
       assert(success0);

       (success0, id) = address(address(0x90)).call(abi.encodeWithSignature("run()", address(this), callArg));   
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

    // function jobExample() pure public returns(uint256){
    //   return 112;
    // }

    function jobExample(address addr, bytes memory id1) pure public returns(uint256){
      return 112;
    }

    function jobExample() pure public returns(uint256){
      return 112;
    }
}