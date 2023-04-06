pragma solidity ^0.5.0;

import "./Base.sol";
import "./Bytes.sol";
import "./U256.sol";
import "./Bool.sol";

contract ByteTest {
    Bytes byteContainer = new Bytes();
    
    function TestBytes() public {     
        require(byteContainer.length() == 0); 
        bytes memory byteArray = new bytes(75);
        for (uint  i = 0; i < 75; i ++) {
            byteArray[i] = 0x41;
        }

        byteContainer.push(byteArray);  
        byteContainer.push(byteArray); 
        require(byteContainer.length() == 2); 

        bytes memory stored = byteContainer.get(1);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        bytes memory elems = new bytes(5);
        for (uint  i = 0; i < elems.length; i ++) {
            elems[i] = 0xaa;
        }
        byteContainer.set(1, elems);
       
        stored = byteContainer.get(0);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        stored = byteContainer.get(1);
        require(stored.length == elems.length); 
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        stored = byteContainer.pop();
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        stored = byteContainer.pop();
        require(byteContainer.length() == 0); 
    }
}