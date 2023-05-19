pragma solidity ^0.5.0;

contract Threading {
    address constant private API = address(0x90); 

    // Append a new task to the queue, the execution only starts when run() is called
    function add(address addr, bytes memory args) public returns(bool, uint256)  {
        (bool success,) = address(API).call(abi.encodeWithSignature("add(address,bytes)", addr, args));   
        return (success, 1);
    }

    // Return the size of the queue
    function length() public returns(uint256)  {
        (, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));   
        return abi.decode(data, (uint256));
    }

    // Remove an item from the queue. Returns true if successful and false otherwise
    function del(uint256 jobID) public returns(bool)  {
        (bool success,) = address(API).call(abi.encodeWithSignature("del(uint256)", jobID));   
        return success;
    }

    // Start processing all the parallel jobs in the queue by specifing the number of threads, the number is between [1, 255]
    function run(uint256 threads) public returns(bool) {
        (bool success,) = address(API).call(abi.encodeWithSignature("run(uint256)", threads));   
        return success;
    } 
    
    // Get the return data by index.
    function get(uint256 id) public returns(bool, bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(uint256)", id));   
        return (success, data);
    }

    // Clear queue and set the length back to zero. Returns true if successful and false otherwise.
    function clear() public returns(bool)  {
        (bool success,) = address(API).call(abi.encodeWithSignature("clear()"));   
        return success;
    }

    //Return True if the queue is empty, False otherwise. 
    function  empty() public returns(bool)  {
        return length() == 0;
    }
}
