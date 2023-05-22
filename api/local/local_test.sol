pragma solidity ^0.5.0;

 import "./U256.sol";

contract Dict {
    U256 keys = new U256();
    mapping(uint256 => uint256) private data ;

    function length() public returns(uint256) { return keys.length();}
    function get(uint256 idx) public returns(uint256)  { return keys.get(idx); }
    function set(uint256 key, uint256 values) public {}
}