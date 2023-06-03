// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

library Storage {
    function purge() public {   
        address(address(0xa0)).call(abi.encodeWithSignature("purge()"));      
    }
} 