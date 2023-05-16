pragma solidity ^0.5.0;

contract Int64Cumulative {
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

    function set(uint256 idx, uint256 v, bool isPositive) public  returns(bool){ // c9 ef ba b9
        (bool success,) = address(API).call(abi.encodeWithSignature("set(bytes, uint256)", idx, v, isPositive));
        return success;
    }
}
