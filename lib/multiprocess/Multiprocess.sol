// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../runtime/Runtime.sol";
import "../bool/Bool.sol";

contract Multiprocess is Base, Revertible {
    uint256 numThreads = 1;
    uint256 serial = 0; 
    constructor (uint256 threads) {
        Base.API = address(0xb0);
        numThreads = threads; 
    } 

    function push(bytes memory elem) public virtual { //9e c6 69 25
        // serial++;
        // Base.setKey(bytes.concat(uuid(), "-", abi.encodePacked(serial)), abi.encode(elem));
        Base.setKey(uuid(), abi.encode(elem));
    }    
 
    function pop() public virtual returns(bytes memory) { // 80 26 32 97
        return abi.decode(Base.popBack(), (bytes));  
    }

    function get(uint256 idx) public virtual  returns(bytes memory)  { // 31 fe 88 d0
        return abi.decode(Base.getIndex(idx), (bytes));  
    }

    function set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        Base.setIndex(idx, abi.encode(elem));   
    }

    function run() public {       
        foreach(abi.encode(numThreads));
    }
}
