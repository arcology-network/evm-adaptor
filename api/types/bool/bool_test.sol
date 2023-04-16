pragma solidity ^0.5.0;

import "./Bool.sol";

contract BoolTest {
    Bool container = new Bool();
    
    constructor() public {     
        require(container.length() == 0); 
    
        container.push(true);
        container.push(false);
        container.push(false);
        container.push(true);
        require(container.length() == 4); 

        require(container.get(0));
        require(!container.get(1));
        require(!container.get(2));
        require(container.get(3));

        container.set(0, false);
        container.set(1, true);
        container.set(2, true);
        container.set(3, false);

        require(!container.get(0));
        require(container.get(1));
        require(container.get(2));
        require(!container.get(3));

        // require(!container.pop());
        // require(container.pop());
        // require(container.pop());
        // require(!container.pop());
        // require(container.length() == 0); 

    }
}