// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Int256Cumulative {
    address constant public API = address(0x86);    
    bytes private id;

    constructor (int256 min, int256 max) public { // 90 54 ce 5f
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new(int256, int256)", min, max));
        id = data; 
        assert(success);
    }
    
    function get() public returns(int256) {  
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("get()", id));
        return abi.decode(data, (int256));
    }

    function add(int256 v) public returns(bool) { // a4 c6 a7 68
        (bool success,) = address(API).call(abi.encodeWithSignature("add(bytes, int256)", id, v));
        return success; 
    }

    function sub(int256 v) public returns(bool) { // c8 da aa ab
        (bool success,) = address(API).call(abi.encodeWithSignature("sub(bytes, int256)", id, v));
        return success;
    }   
}
