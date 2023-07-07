// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Bytes32.sol";

contract Bytes32Test {
    Bytes32 container = new Bytes32();
    
    constructor() {     
        require(container.length() == 0); 
   
        bytes32 hash0 = keccak256(abi.encodePacked("0"));
        bytes32 hash1 = keccak256(abi.encodePacked("1"));
        bytes32 hash2 = keccak256(abi.encodePacked("2"));
        bytes32 hash3 = keccak256(abi.encodePacked("3"));

        container.push(hash0);
        container.push(hash1);
        container.push(hash2);
        container.push(hash3);
        require(container.length() == 4); 

        require(container.get(0) == hash0);
        require(container.get(1) == hash1);
        require(container.get(2) == hash2);
        require(container.get(3) == hash3);

        container.set(0, hash3);
        container.set(1, hash2);
        container.set(2, hash1);
        container.set(3, hash0);

        require(container.get(0) == hash3);
        require(container.get(1) == hash2);
        require(container.get(2) == hash1);
        require(container.get(3) == hash0);

        require(container.pop() == hash0);
        require(container.pop() == hash1);
        require(container.pop() == hash2);

        container.clear();
        require(container.length() == 0);       
    }
}