// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

contract Bytes32 {
    Base base;

    event logMsg(string message);

    constructor() {  base = new Base(address(0x84), false); }
    function length() public returns(uint256) { return base.length();}

    function pop() public returns(bytes32) { 
        return abi.decode(base.pop(), (bytes32));  
    }

    function push(bytes32 elem) public { 
       base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }   

    function get(uint256 idx) public returns(bytes32)  { 
        return abi.decode(base.get(idx), (bytes32));
    }

    function set(uint256 idx, bytes32 elem) public {
        base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
    }

    function clear() public { // 7a fa 62 38
        base.clear();
    }
}


// contract Bytes32 is Base {
//     constructor(address addr, bool local) Base(addr, local) {}

//     function push(bytes32 elem) public virtual{ //9e c6 69 25
//         Base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
//     }    

//     function pop() public virtual override returns(bytes memory) { // 80 26 32 97
//         return abi.decode(Base.pop(), (bytes));  
//     }

//     function get(uint256 idx) public virtual override returns(bytes memory)  { // 31 fe 88 d0
//         return abi.decode(Base.get(idx), (bytes));  
//     }

//     function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
//         Base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));     
//     }
// }
