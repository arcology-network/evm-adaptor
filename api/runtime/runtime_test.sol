// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Runtime.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";


contract Deployee is Localizer{ 
    uint256[2] num;    
    constructor() {
        num[0] = 11;
        num[1] = 12;
        require(
            num[0] == 11 && 
            num[1] == 12        
        );
    }    

    function check() public {
        require(num[0] == 111);
    }
}

contract Deployer { 
   constructor() {
       Deployee deployee = new Deployee();
    //    deployee.check();
  }

  function deployers() public {
       Deployee deployee = new Deployee();
  }
}
