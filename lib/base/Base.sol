// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;
import "../runtime/Runtime.sol";

contract Base is Runtime{
    address internal API = address(0x84);    

    constructor () {
        (bool success,) = address(API).call(abi.encodeWithSignature("new()", true));       
        require(success);
    }

    function length() public returns(uint256) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));
        require(success);
        return abi.decode(data, (uint256));
    }

    // The initial length of the container at the current block height
    function peek() public returns(bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("peek()"));
        require(success);
        return data;  
    } 

    function popBack() public virtual returns(bytes memory) { // 80 26 32 97
        bytes memory v = getIndex(length() - 1);
        delIndex(length() - 1);
        return v;
    }

    function setIndex(uint256 idx, bytes memory encoded) public { // 7a fa 62 38
        address(API).call(abi.encodeWithSignature("setIndex(uint256,bytes)", idx, encoded));     
    }

    function setKey(bytes memory key, bytes memory elem) public {
        address(API).call(abi.encodeWithSignature("setKey(bytes,bytes)", key, elem));
    }

    function delIndex(uint256 idx) public { // 7a fa 62 38
        address(API).call(abi.encodeWithSignature("delIndex(uint256)", idx));     
    }

    function delKey(bytes memory key) public {
        address(API).call(abi.encodeWithSignature("delKey(bytes)", key));
    }

    function getIndex(uint256 idx) public virtual returns(bytes memory) { // 31 fe 88 d0
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getIndex(uint256)", idx));
        if (success) {
            return abi.decode(data, (bytes));  
        }
        return data;
    }

    function getKey(bytes memory key) public returns(bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getKey(bytes)", key));
        if (success) {
            return abi.decode(data, (bytes));  
        }
        return data;
    }

    // Clear the data
    function clear() public {
        address(API).call(abi.encodeWithSignature("clear()"));       
    }

    function foreach(bytes memory data) public {
        address(API).call(abi.encodeWithSignature("foreach(bytes)", data));       
    }
}
