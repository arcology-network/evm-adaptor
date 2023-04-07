pragma solidity ^0.5.0;

import "./Base.sol";

contract Address {
    Base base;

    event logMsg(string message);

    constructor  () public {
        base = new Base();
    }

    function length() public returns(uint256) {  // 58 94 13 33
        return base.length();
    }

    function pop() public returns(address) { // 80 26 32 97
        return abi.decode(base.pop(), (address ));  
    }

    function push(bytes memory elem) public { //9e c6 69 25
        base.push(abi.encodeWithSignature("push(bytes, bytes)",  base.id(), elem));
    }   

    function get(uint256 idx) public returns(address)  { // 31 fe 88 d0
        return abi.decode(base.get(idx), (address));  
    }

    function set(uint256 idx, address elem) public { // 7a fa 62 38
        base.set(abi.encodeWithSignature("set(bytes, uint256, bytes)", base.id(), idx, elem));     
    }
}

