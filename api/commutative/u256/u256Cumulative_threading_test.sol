pragma solidity ^0.5.0;

import "./U256Cumulative.sol";
import "./Threading.sol";

contract ThreadingCumulativeU256 {
    U256Cumulative cumulative = new U256Cumulative(0, 100); 
    
    function call() public {
       Threading mp = new Threading();
       mp.add(address(this), abi.encodeWithSignature("add(uint256)", 1));
       mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));      
       mp.run(1);

      assert(cumulative.get() == 3);
    }

    function add(uint256 elem) public { //9e c6 69 25
       cumulative.add(elem);
    }  
}
