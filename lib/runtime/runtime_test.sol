// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Runtime.sol";
import "../bool/Bool.sol";
import "../u256cum/U256Cum.sol";


contract ResettableDeployer {
    // uint256 public immutable a;
    // This is a valid constructor
    TestResettable resettable;
    constructor () {
        resettable = new TestResettable();
        resettable.set();
    }
    
    // This is invalid and will not compile
    function afterCheck () public {
        resettable.afterCheck();
    }
}

contract TestResettable is Revertible {
    uint256 [2] public array; 
    constructor () {
        array[0] = 10;
        array[1] = 11;     
    }

    function set() public {
        require(array[0] == 10);
        require(array[1] == 11);   
    }

    function afterCheck() public {
        require(array[0] == 10);
        require(array[1] == 11);   
        rollback();  

        require(array[0] == 10);
        require(array[1] == 11);   
        
        array[0] = 100;
        array[1] = 111;
   
        require(array[0] == 100);
        require(array[1] == 111);
        rollback();  

        require(array[0] == 10);
        require(array[1] == 11);  
    }
}





// contract TestImmutable {
//     // uint256 public immutable a;
//     // This is a valid constructor
//     constructor () {
//         address(0xa0).call(abi.encodeWithSignature("Reset()"));
//     }
    
//     // This is invalid and will not compile
//     function afterCheck () public {
//         address(0xa0).call(abi.encodeWithSignature("Reset()"));
//     }
// }