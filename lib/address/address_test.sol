// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Address.sol";

contract AddressTest {
    Address container = new Address();
    
    constructor() {     
        require(container.length() == 0);    
        address addr1 = 0x1111111110123456789012345678901234567890;
        address addr2 = 0x2222222220123456789012345678901234567890;
        address addr3 = 0x3333337890123456789012345678901234567890;
        address addr4 = 0x4444444890123456789012345678901234567890;

        container.push(addr1);
        container.push(addr2);
        require(container.length() == 2); 

        require(container.get(0) == addr1);
        require(container.get(1) == addr2);
 
        container.set(0, addr3);
        container.set(1, addr4);

        require(container.get(0) == addr3);
        require(container.get(1) == addr4);

        // require(container.pop() == addr4);
        // require(container.pop() == addr3);
        // require(container.length() == 0); 
    } 
} 