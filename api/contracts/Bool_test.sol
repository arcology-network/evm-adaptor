pragma solidity ^0.5.0;

import "./Bool.sol";

contract BoolTest {
    Bool boolContainer = new Bool();
    
    constructor() public {     
        require(boolContainer.length() == 0); 
    
        boolContainer.push(true);
        boolContainer.push(false);
        boolContainer.push(false);
        boolContainer.push(true);
        require(boolContainer.length() == 4); 

        require(boolContainer.get(0));
        require(!boolContainer.get(0));
        require(!boolContainer.get(0));
        require(boolContainer.get(0));

        boolContainer.set(0, false);
        boolContainer.set(1, true);
        boolContainer.set(2, true);
        boolContainer.set(3, false);

        require(!boolContainer.get(0));
        require(boolContainer.get(0));
        require(boolContainer.get(0));
        require(!boolContainer.get(0));

        require(!boolContainer.pop());
        require(boolContainer.pop());
        require(boolContainer.pop());
        require(!boolContainer.pop());
    }
}