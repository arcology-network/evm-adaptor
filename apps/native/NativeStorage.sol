// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract NativeStorage {   
        uint256 x = 1 ;
        uint256 y = 100 ;

   constructor() { 
        incrementX();
        incrementY();
        require(x == 2);
        require(y == 102);
   }

    function incrementX() public {
        x ++;
    }

    function incrementY() public {
       y += 2;
    }

    function checkX(uint256 value) view public {
        require(x == value);
    }

    function checkY(uint256 value) view public {
         require(y == value);
    }
}