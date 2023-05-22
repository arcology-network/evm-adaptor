pragma solidity ^0.5.0;

library Storage {
    function forget() public {   
        address(address(0xa0)).call(abi.encodeWithSignature("forget()"));      
    }
} 