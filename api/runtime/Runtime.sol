// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {
    address constant API = address(0xa0) ;    
    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(API).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}

contract Resettable is Runtime { 
    function reset() public {
        address(API).call(abi.encodeWithSignature("Reset()"));     
    }
}