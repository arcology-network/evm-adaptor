pragma solidity ^0.5.0;

import "./Multiprocess.sol";

contract MultiprocessTest {
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

    function isPrime(uint256 n) public pure returns(bool) {
        // Check if n is less than 2.
        if (n < 2) {
            return false;
        }
        // Check if n is 2 or 3.
        if (n == 2 || n == 3) {
            return true;
        }
        // Check if n is divisible by 2 or 3.
        if (n % 2 == 0 || n % 3 == 0) {
            return false;
        }
        // Check for other divisors up to sqrt(n).
        for (uint256 i = 5; i*i <= n; i += 6) {
            if (n % i == 0 || n % (i + 2) == 0) {
                return false;
            }
        }
        return true;
    }
}