pragma solidity ^0.5.0;

contract Parallel {
    address constant public API = address(0x90); 

    function newJob (address addr, bytes memory args) public returns(bool, uint32)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("newJob(address, uint32), returns(bool, uint32)", addr, args));   
        return (success, abi.decode(data, (uint32)));
    }

    function job (uint32 id) public returns(bool, bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("job(uint32), returns(bool, uint32)", id));   
        return (success, abi.decode(data, (bytes)));
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
