pragma solidity ^0.5.0;

import "./Base.sol";

contract Address {
    Base base;

    event logMsg(string message);

    constructor  () public {  base = new Base(); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(address) { 
        return abi.decode(abi.decode(base.pop(), (bytes)), (address));  
    }

    function push(address elem) public { 
       base.push(abi.encodeWithSignature("push(bytes,bytes)",  base.id(), abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(address)  { 
        return abi.decode(abi.decode(base.get(idx), (bytes)), (address));
    }

    function set(uint256 idx, address elem) public {
        base.set(abi.encodeWithSignature("set(bytes,uint256,bytes)", base.id(), idx, abi.encode(elem)));        
    }
}

