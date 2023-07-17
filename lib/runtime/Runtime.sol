// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {    
    event logMsg(string message);

    function pid() public returns(bytes memory args) {
        (,bytes memory randome) = address(0xa0).call(abi.encodeWithSignature("pid()"));     
        return randome;
    }

    function uuid() public returns(bytes memory) {
        (,bytes memory id) = address(0xa0).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}

contract Revertible { 
    function rollback() public {
        address(0xa0).call(abi.encodeWithSignature("Reset()"));     
    }
}