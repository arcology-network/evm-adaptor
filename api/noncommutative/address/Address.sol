// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;


import "../base/Base.sol";

// contract Address {
//     Base base;

//     event logMsg(string message);

//     constructor() {  base = new Base(address(0x84), false); }
//     function length() public returns(uint256) { return base.length();}

//     function pop() public returns(address) { 
//         return abi.decode(Base.popBack(), (address));  
//     }

//     function push(address elem) public { 
//        Base.pushBack(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
//     }   

//     function get(uint256 idx) public returns(address)  { 
//         return abi.decode(Base.getElem(idx), (address));
//     }

//     function set(uint256 idx, address elem) public {
//         Base.setElem(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));        
//     }

//     function clear() public { // 7a fa 62 38
//         base.clear();
//     }
// }

contract Address is Base {
    constructor() Base(address(0x84), false) {}

    function push(address elem) public virtual{ //9e c6 69 25
        Base.pushBack(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }    

    function pop() public virtual returns(address) { // 80 26 32 97
        return abi.decode(Base.popBack(), (address));  
    }

    function get(uint256 idx) public virtual returns(address)  { // 31 fe 88 d0
        return abi.decode(Base.getElem(idx), (address));  
    }

    function set(uint256 idx, address elem) public { // 7a fa 62 38
        Base.setElem(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));     
    }
}
