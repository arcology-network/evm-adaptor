// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Int256Cumulative.sol";
import "../../threading/Threading.sol";

contract ThreadingInt64 {
    Int256Cumulative container = new Int256Cumulative(0, 100);

   //  function call() public {
   //     Threading mp = new Threading(1);
   //     mp.add(address(this), abi.encodeWithSignature("add(int64)", 1));
   //     mp.add(address(this), abi.encodeWithSignature("add(int64)", 2));      
       


   //     mp.run();
   //     assert(container.length() == 1 ); 
   //  }

   //  function push(int64 elem) public { //9e c6 69 25
   //     container.push(elem);
   //  }  
}
