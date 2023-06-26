// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Bytes32N.sol";

contract Bytes32NTest {
    Bytes32N container;

    constructor() {    
        bytes32 arr1 = keccak256(abi.encodePacked("0"));
        bytes32 arr2 = keccak256(abi.encodePacked("1"));
        bytes32 arr3 = keccak256(abi.encodePacked("2"));
        bytes32 arr4 = keccak256(abi.encodePacked("3"));

        container = new Bytes32N(4, arr1);
        require(container.length() == 4); 
        
        require((container.get(0)) == (arr1));
        require((container.get(1)) == (arr1));
        require((container.get(2)) == (arr1));
        require((container.get(3)) == (arr1));

        container.set(0, arr4);
        container.set(1, arr3);
        container.set(2, arr2);
        container.set(3, arr1);

        require((container.get(0)) == (arr4));
        require((container.get(1)) == (arr3));
        require((container.get(2)) == (arr2));
        require((container.get(3)) == (arr1));    
    }
}