// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Runtime.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";


contract Deployee is Localizer{ 
    uint256[2] public num ;    
    constructor() {
        num[0] = 11;
        num[1] = 12;
        require(
            num[0] == 11 && 
            num[1] == 12        
        );
    }    

    function check() public returns(bool){
        return (num[0] == 11);
    }
}

contract NonLocalizedDeployee { 
    uint256[2] public num ;    
    constructor() {
        num[0] = 11;
        num[1] = 12;
        require(
            num[0] == 11 && 
            num[1] == 12        
        );
    }    

    function check() public returns(bool){
        return (num[0] == 11);
    }
}


contract Deployer { 
   Deployee localizedDeployee;
   NonLocalizedDeployee nonLocalizedDeployee; 

   constructor() {
       localizedDeployee = new Deployee();
       nonLocalizedDeployee = new NonLocalizedDeployee();

       require(localizedDeployee.check());
       require(nonLocalizedDeployee.check());    
  }

  function afterCheck() public {
    // localizedDeployee.check();
    require(localizedDeployee.check());
    require(nonLocalizedDeployee.check());
  }
}
