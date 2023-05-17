pragma solidity ^0.5.0;

contract Int256Cumulative {
    address constant public API = address(0x85);    
    bytes private id;

    constructor (int256 min, int256 max) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("New(uint256, uint256, uint256)", min, max));
        id = data; 
        assert(success);
    }
    
    function get() public returns(uint256) {  
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("get()", id));
        return abi.decode(data, (uint256));
    }

    function add(int256 v) public returns(bool) { 
        (bool success,) = address(API).call(abi.encodeWithSignature("add(bytes, int256)", id, v));
        return success; 
    }

    function sub(int256 v) public returns(bool) { 
        (bool success,) = address(API).call(abi.encodeWithSignature("sub(bytes, int256)", id, v));
        return success;
    }   
}
