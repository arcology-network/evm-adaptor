pragma solidity ^0.5.0;

import "./Int64.sol";

contract Int64Fixed {
    Int64 array = new Int64();
    constructor  (uint length, int64 initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(int64)  {return array.get(idx);}
    function set(uint256 idx, int64 elem) public { array.set(idx, elem); }
}
