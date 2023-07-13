// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

contract Runtime {
    // address constant private API = address(0xa0); 
    constructor (string memory property) {
        address(0xa0).call(abi.encodeWithSignature("set(string)", property));     
    }

    function uuid() public returns(bytes memory args) {
        (,bytes memory id) = address(0xa0).call(abi.encodeWithSignature("uuid()"));     
        return id;
    }
}

contract Localizer is Runtime{ 
    constructor() Runtime("local=true"){}
}


// library runtimelib {
//     address constant private API = address(0xa0); 
//     function set(string calldata property) public returns(bool) {
//         (bool ok,) = address(API).call(abi.encodeWithSignature("set(string)", property));     
//         return ok;
//     }
// }