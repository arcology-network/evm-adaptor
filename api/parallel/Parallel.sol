pragma solidity ^0.5.0;

contract Parallel {
    address constant public API = address(0x90); 

    function newJob (bytes memory funcCall) public returns(bool, uint32)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("NewJob(bytes), returns(bool, uint32)", funcCall));   
        return (success, abi.decode(data, (uint32)));
    }

    function delJob (uint32 jobID) public returns(bool)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("DelJob(uint32), returns(bool)", jobID));   
        return success;
    }

    function run() public returns(bool, bool[] memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("run(), returns(bool, bool[])"));   
        return (success, abi.decode(data, (bool[])));
    } 

    function clear() public returns(bool)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("clear(), returns(bool)"));   
        return success;
    }
}
