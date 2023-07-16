// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {    
    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(0xa0).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}

contract Revertible { 
    function rollback() public {
        address(0xa0).call(abi.encodeWithSignature("Reset()"));     
    }
}