pragma solidity ^0.5.0;

import "./Bytes32.sol";

contract Bytes32Fixed {
    Bytes32 array = new Bytes32();
    constructor  (uint length, bytes32 initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(bytes32)  {return array.get(idx);}
    function set(uint256 idx, bytes32 elem) public { array.set(idx, elem); }
}
