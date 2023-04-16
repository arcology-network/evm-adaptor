pragma solidity ^0.5.0;

contract Parallel {
    address constant private API = address(0x90); 

    function addJob (address addr, bytes memory args) public returns(bool, uint32)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("addJob(address,bytes)", addr, args));   
        return (success, 1);
    }

    function job (uint32 id) public returns(bool, bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("job(uint)", id));   
        return (success, abi.decode(data, (bytes)));
    }

    function delJob (uint32 jobID) public returns(bool)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("delJob(uint)", jobID));   
        return success;
    }

    function run() public returns(bool, bool[] memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("run()"));   
        return (success, abi.decode(data, (bool[])));
    } 

    function clear() public returns(bool)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("clear()"));   
        return success;
    }

    function export (bytes memory data) public returns(bool, uint32)  {
         (bool success, bytes memory data) = address(API).call(data); 
        return (true, 1);
    }

    function jobExample(bytes memory id, bytes memory id2) pure public returns(bytes memory){
      return id2;
    }

}