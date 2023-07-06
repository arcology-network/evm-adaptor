// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

address constant PARALLEL_API = address(0xb0);   
address constant BASE_API= address(0x84);   

contract Foo {
    // Doesn't need to declare MY_CONSTANT as this is done
    // at the file level in constants.sol

    // Use the imported definition from ./constants.sol
    function myFunction() external pure returns (uint256) {
        return 12; 
    }
}