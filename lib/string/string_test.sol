// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./String.sol";

contract StringTest {
    String container = new String();
    
    constructor() {     
        require(container.length() == 0); 
   
        string memory str0 = "Test string 0";
        string memory str1 = "Test string 1";
        string memory str2 = "Test string 2";
        string memory str3 = "Test string 3";

        container.push(str0);
        container.push(str1);
        container.push(str2);
        container.push(str3);
        require(container.length() == 4); 

        require(keccak256(bytes(container.get(0))) == keccak256(bytes(str0)));
        require(keccak256(bytes(container.get(1))) == keccak256(bytes(str1)));
        require(keccak256(bytes(container.get(2))) == keccak256(bytes(str2)));
        require(keccak256(bytes(container.get(3))) == keccak256(bytes(str3)));

        container.set(0, str3);
        container.set(1, str2);
        container.set(2, str1);
        container.set(3, str0);

        require(keccak256(bytes(container.get(3))) == keccak256(bytes(str0)));
        require(keccak256(bytes(container.get(2))) == keccak256(bytes(str1)));
        require(keccak256(bytes(container.get(1))) == keccak256(bytes(str2)));
        require(keccak256(bytes(container.get(0))) == keccak256(bytes(str3)));

        require(keccak256(bytes(container.pop())) == keccak256(bytes(str0)));
        require(keccak256(bytes(container.pop())) == keccak256(bytes(str1)));
        require(keccak256(bytes(container.pop())) == keccak256(bytes(str2)));
        require(keccak256(bytes(container.pop())) == keccak256(bytes(str3)));

        require(container.length() == 0);       
    }
}