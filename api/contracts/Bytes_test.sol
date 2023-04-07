pragma solidity ^0.5.0;

import "./Bytes.sol";

contract ByteTest {
    Bytes container = new Bytes();
    
    constructor() public {     
        require(container.length() == 0); 
        bytes memory byteArray = new bytes(75);
        for (uint  i = 0; i < 75; i ++) {
            byteArray[i] = 0x41;
        }

        container.push(byteArray);  
        container.push(byteArray); 
        require(container.length() == 2); 

        bytes memory stored = container.get(1);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        bytes memory elems = new bytes(5);
        for (uint  i = 0; i < elems.length; i ++) {
            elems[i] = 0xaa;
        }
        container.set(1, elems);
       
        stored = container.get(0);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        stored = container.get(1);
        require(stored.length == elems.length); 
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        stored = container.pop();
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        stored = container.pop();
        require(container.length() == 0); 
    }
}