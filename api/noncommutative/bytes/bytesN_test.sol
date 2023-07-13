// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./BytesN.sol";

contract BytesNTest {
    BytesN container;

    constructor() {    
        bytes memory arr1 = '0x1000000000000000000000000000000000000000000000000000000000000001';
        bytes memory arr2 = '0x2000000000000000000000000000000000000000000000000000000000000002';
        bytes memory arr3 = '0x3000000000000000000000000000000000000000000000000000000000000003';
        bytes memory arr4 = '0x4000000000000000000000000000000000000000000000000000000000000004';

        container = new BytesN(4, arr1);
        require(container.length() == 4); 
        
        require(keccak256(container.get(0)) == keccak256(arr1));
        require(keccak256(container.get(1)) == keccak256(arr1));
        require(keccak256(container.get(2)) == keccak256(arr1));
        require(keccak256(container.get(3)) == keccak256(arr1));

        container.set(0, arr4);
        container.set(1, arr3);
        container.set(2, arr2);
        container.set(3, arr1);

        require(keccak256(container.get(0)) == keccak256(arr4));
        require(keccak256(container.get(1)) == keccak256(arr3));
        require(keccak256(container.get(2)) == keccak256(arr2));
        require(keccak256(container.get(3)) == keccak256(arr1));    
    }
}