pragma solidity ^0.5.0;

import "./Dynamic.sol";

contract U256Fixed {
    U256 array;
    constructor  (uint length, uint256 value) public {  
        array = new U256(); 
        for (uint i = 0; i < length; i ++) {
            array.push(value);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(uint256)  {return array.get(idx);}
    function set(uint256 idx, uint256 elem) public { array.set(idx, elem); }
}
