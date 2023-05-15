pragma solidity ^0.5.0;

import "./Base.sol";

contract Int64 {
    Base base;

    constructor  () public {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(int) { // 80 26 32 97
        return abi.decode(abi.decode(base.pop(), (bytes)), (int));  
    }

    function push(int64 elem) public{ //9e c6 69 25
       base.push(abi.encodeWithSignature("push(bytes,bytes)",  base.id(), abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(int64)  { // 31 fe 88 d0
        return abi.decode(abi.decode(base.get(idx), (bytes)), (int64));
    }

    function set(uint256 idx, int64 elem) public { // 7a fa 62 38
        base.set(abi.encodeWithSignature("set(bytes,uint256,bytes)", base.id(), idx, abi.encode(elem)));        
    }
}
