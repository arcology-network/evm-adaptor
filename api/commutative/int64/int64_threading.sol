pragma solidity ^0.5.0;

import "./Int64.sol";
import "./Threading.sol";

contract ThreadingInt64 {
    Int64 container = new Int64();

    function call() public {
       Threading mp = new Threading();
       mp.add(address(this), abi.encodeWithSignature("push(int64)", 1));
       mp.add(address(this), abi.encodeWithSignature("push(int64)", 2));      
       assert(container.length() == 0 && mp.length() == 2); 
       
       mp.run(1);
       assert(container.length() == 1 ); 
    }

    function push(int64 elem) public { //9e c6 69 25
       container.push(elem);
    }  
}
