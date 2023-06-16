// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract U256Cumulative {
    address constant public API = address(0x85);    
    bytes private id;

    constructor (uint256 min, uint256 max) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("New(uint256, uint256, uint256)", min, max));
        id = data; 
        assert(success);
    }
    
    function get() public returns(uint256) {  
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("get()", id));
        return abi.decode(data, (uint256));
    }

    function add(uint256 v) public returns(bool) { 
        (bool success,) = address(API).call(abi.encodeWithSignature("add(bytes, uint256)", id, v));
        return success; 
    }

    function sub(uint256 v) public returns(bool) { 
        (bool success,) = address(API).call(abi.encodeWithSignature("sub(bytes, uint256)", id, v));
        return success;
    }   

    function min() public returns(uint256) { 
        (, bytes memory data) = address(API).call(abi.encodeWithSignature("min(bytes)", id));
        return abi.decode(data, (uint256));
    }  

    function max() public returns(uint256) { 
        (, bytes memory data) = address(API).call(abi.encodeWithSignature("max(bytes)", id));
        return abi.decode(data, (uint256));
    }    
}
