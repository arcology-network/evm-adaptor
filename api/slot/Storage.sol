pragma solidity ^0.5.0;

library Storage {
    function purge() public {   
        address(address(0xa0)).call(abi.encodeWithSignature("purge()"));      
    }
} 