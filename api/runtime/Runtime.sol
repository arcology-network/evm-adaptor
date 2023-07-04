// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {
    address constant private API = address(0xa0); 

    function deferred(uint256 gaslimit, address addr, bytes memory args) public returns(bool, bytes memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("defer(uint256,bytes)", gaslimit, addr, args));     
        if (!success) {
            (success, data) = addr.call(args);
            return (success, data);
        }
        return (success,data);
    }


    function singleton() public returns(bool) {
        (bool ok,) = address(API).call(abi.encodeWithSignature("singleton()"));     
        return ok;
    }

    function localize(uint256 slot) public returns(bool) {
        (bool ok,) = address(API).call(abi.encodeWithSignature("singleton(uint256)", slot));     
        return ok;
    }

    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(API).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}
