// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {
    constructor (bytes memory property) {      // address constant private API = address(0xa0); 
        address(0xa0).call(property);     
    }

    function exists(bytes memory key) public returns(bool) {
        (,bytes memory id) = address(0xa0).call(abi.encodeWithSignature("exists(bytes)", key));
        return abi.decode(id, (bool));
    }

    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(0xa0).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}

contract Revertible { 
    function revert() public {
        address(0xa0).call(abi.encodeWithSignature("Reset()"));     
    }
}