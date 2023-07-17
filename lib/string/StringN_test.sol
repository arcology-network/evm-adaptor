// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./StringN.sol";

contract Int64NTest {
    StringN container;

    constructor() {    
        string memory _1 = "1111";
        string memory _2 = "2222";
        string memory _3 = "3333";
        string memory _4 = "4444";

        container = new StringN(4, _1);
        require(container.length() == 4); 
        
        require(keccak256(bytes(container.get(0))) == keccak256(bytes(_1)));
        require(keccak256(bytes(container.get(1))) == keccak256(bytes(_1)));
        require(keccak256(bytes(container.get(2))) == keccak256(bytes(_1)));
        require(keccak256(bytes(container.get(3))) == keccak256(bytes(_1)));

        container.set(0, _4);
        container.set(1, _3);
        container.set(2, _2);
        container.set(3, _1);

        require(keccak256(bytes(container.get(0))) == keccak256(bytes(_4)));
        require(keccak256(bytes(container.get(1))) == keccak256(bytes(_3)));
        require(keccak256(bytes(container.get(2))) == keccak256(bytes(_2)));
        require(keccak256(bytes(container.get(3))) == keccak256(bytes(_1)));
    }
}