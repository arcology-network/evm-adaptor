pragma solidity ^0.5.0;

contract U256 {
    address constant public API = address(0x85);    
    bytes private id;

    constructor (uint256 v, uint256 min, uint256 max) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("New(uint256, uint256, uint256) returns(bytes)", v, min, max));
        id = data; 
        assert(success);
    }
    
    function get() public returns(uint256) {  
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get() returns(uint256, bool)", id));
        return abi.decode(data, (uint256));
    }

    function add(uint256 v) public returns(bool) { 
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("add(bytes, uint256) returns(bool)", id, v));
        return success; 
    }

    function sub(uint256 v) public returns(bool) { 
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("sub(bytes, uint256) returns(bool)", id, v));
        return success;
    }   
}
