pragma solidity ^0.5.0;

contract U64 {
    address constant public API = address(0x85);    
    bytes private id;

    constructor (uint64 v) public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("New() returns(bytes)"));
        require(success);
        id = data; 
    }

    function add(uint64 v) public returns(bool) { // 80 26 32 97
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("add(bytes, uint64) returns(bool)", id, v));
        return success; 
    }
    
    function get() public returns(uint64, bool) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(bytes) returns(bool, uint64)", id));
        return (abi.decode(data, (uint64)), success);
    }
}
