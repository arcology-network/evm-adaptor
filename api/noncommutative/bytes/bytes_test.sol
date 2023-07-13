// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

// import "./Bytes.sol";

import "../base/Base.sol";

// contract Bytes is Base {
//     constructor(address addr, bool local) Base(addr, local) {}

//     function push(bytes memory elem) public virtual override { //9e c6 69 25
//         Base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem))); 
//     }    

//     function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
//         Base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));     
//     }
// }

contract Bytes is Base {
    constructor(address addr, bool local) Base(addr, local) {}

    function push(bytes memory elem) public virtual override { //9e c6 69 25
        Base.push(abi.encodeWithSignature("push(bytes)", abi.encode(elem)));
    }    

    function pop() public virtual override returns(bytes memory) { // 80 26 32 97
        return abi.decode(Base.pop(), (bytes));  
    }

    function get(uint256 idx) public virtual override returns(bytes memory)  { // 31 fe 88 d0
        return abi.decode(Base.get(idx), (bytes));  
    }

    function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        Base.set(abi.encodeWithSignature("set(uint256,bytes)", idx, abi.encode(elem)));     
    }
}

contract ByteTest {
    Bytes container = new Bytes(address(0x84), false);    

    constructor() {       
        require(container.length() == 0); 
 
        bytes memory arr1 = '0x1000000000000000000000000000000000000000000000000000000000000001';
        bytes memory arr2 = '0x0x2000000000000000000000000000000000000000000000000000000000000002';

        container.push(arr1);  
        container.push(arr2); 

        require(container.length() == 2); 
        require(keccak256(container.get(0)) == keccak256(arr1));
        require(keccak256(container.get(1)) == keccak256(arr2));

        container.set(1, arr2);       
        require(keccak256(container.get(0)) == keccak256(arr1));
        require(keccak256(container.get(1)) == keccak256(arr2));
        require(keccak256(container.pop()) == keccak256(arr2));

        container.clear();
        require(container.length() == 0); 
    }
}