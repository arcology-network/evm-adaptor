// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Concurrency {
    address constant private API = address(0xa0); 

    function deferred(address addr, bytes memory args) public returns(bool, bytes memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("deferred()", addr, args));     
        if (!success) {
            (success, data) = addr.call(args);
            return (success, data);
        }
        return (success,data);
    }

    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(API).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}
