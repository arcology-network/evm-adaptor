// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Bool.sol";

contract BoolTest {
    Bool boolContainer = new Bool();
    
    constructor() {     
        require(boolContainer.length() == 0); 
    
        boolContainer.push(true);
        boolContainer.push(false);
        boolContainer.push(false);
        boolContainer.push(true);
        require(boolContainer.length() == 4); 

       require(boolContainer.get(0));
        require(!boolContainer.get(1));
        require(!boolContainer.get(2));
        require(boolContainer.get(3));

        boolContainer.set(0, false);
        boolContainer.set(1, true);
        boolContainer.set(2, true);
        boolContainer.set(3, false);

        require(!boolContainer.get(0));
        require(boolContainer.get(1));
        require(boolContainer.get(2));
        require(!boolContainer.get(3));

        require(!boolContainer.pop());
        require(boolContainer.pop());
        require(boolContainer.pop());
        require(!boolContainer.pop());
        require(boolContainer.length() == 0);  
    }
}