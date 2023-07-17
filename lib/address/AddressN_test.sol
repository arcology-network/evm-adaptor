// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./AddressN.sol";

contract U256NTest {
    AddressN container;

    constructor() {    
        address addr1 = 0x1111111110123456789012345678901234567890;
        address addr2 = 0x2222222220123456789012345678901234567890;
        address addr3 = 0x3333337890123456789012345678901234567890;
        address addr4 = 0x4444444890123456789012345678901234567890;

        container = new AddressN(4, addr1);
        require(container.length() == 4); 

        require(container.get(0) == addr1);
        require(container.get(1) == addr1);
        require(container.get(2) == addr1);
        require(container.get(3) == addr1);

        container.set(0, addr4);
        container.set(1, addr3);
        container.set(2, addr2);
        container.set(3, addr1);

        require(container.get(0) == addr4);
        require(container.get(1) == addr3);
        require(container.get(2) == addr2);
        require(container.get(3) == addr1);    
    }
}