pragma solidity ^0.5.0;

import "./Bytes.sol";

contract BytesN {
    Bytes array = new Bytes();
    constructor  (uint length, bytes memory initialV) public {  
        for (uint i = 0; i < length; i ++) {
            array.push(initialV);
        }
    }

    function length() public returns(uint256) { return array.length();}
    function get(uint256 idx) public returns(bytes memory)  {return array.get(idx);}
    function set(uint256 idx, bytes memory elem) public { array.set(idx, elem); }
}
