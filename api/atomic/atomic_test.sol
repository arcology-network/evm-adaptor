pragma solidity ^0.5.0;

import "./Atomic.sol";
import "./Threading.sol";
import "./Bool.sol";
import "./U256Cumulative.sol";

 // this should fail because the deferred call hasn't been processed yet. 
contract AtomicDeferredInThreadingTest {  
     bytes32[2] results;
     Atomic atomic = new Atomic(); 
     function call() public  { 
        bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
  
        Threading mp = new Threading(2);
        mp.add(800000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 0,data));
        mp.add(800000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 1,data));
        require(mp.length() == 2);
        mp.run();
       
        mp.clear();
        assert(mp.length() == 0);      
    }

    function hasher(uint256 idx, bytes memory data) public {
       results[idx] = keccak256(data);
       atomic.deferred(300000, address(this), abi.encodeWithSignature("example()"));
    }

    function example() public {
       results[0] = keccak256("0");
       results[1] = keccak256("1");
    }

    function PostCheck() public view{
       assert(results[0] == keccak256("0"));
       assert(results[1] == keccak256("1"));
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

contract AtomicDeferredBoolContainerTest {  
     Bool container = new Bool();
     bytes32[2] results;
     Atomic atomic = new Atomic(); 
     function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
   
       Threading mp = new Threading(2);
       mp.add(1000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 0,data));
       mp.add(2000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 1,data));
       mp.add(2000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 2,data));
       require(mp.length() == 3);
       require(container.length() == 0);
       mp.run();      
       require(container.length() == 6);
       
       assert(results[0] == keccak256("0"));
       assert(results[1] == keccak256("1"));

       container.push(false);
       container.push(false);

       require(container.length() == 8);
       mp.clear();
       assert(mp.length() == 0);      
    }

    function hasher(uint256 idx, bytes memory data) public {
       container.push(true);
       results[idx] = keccak256(data);
       atomic.deferred(500000, address(this), abi.encodeWithSignature("example()"));
    }

    function example() public {
       container.push(true);
       container.push(true);
       container.push(true);
       results[0] = keccak256("0");
       results[1] = keccak256("1");
    }

    function PostCheck() public view {
        assert(results[0] == keccak256("0"));
        assert(results[1] == keccak256("1"));
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

contract AtomicMultiDeferredWithBoolContainerTest {  
     Bool container = new Bool();
     U256Cumulative u256Comulative = new U256Cumulative(0, 100);

     bytes32[2] results;
     Atomic atomic = new Atomic(); 
     function call() public  { 
        bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
   
        Threading mp = new Threading(2);
        mp.add(1000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 0,data));
        mp.add(2000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 1,data));
        mp.add(2000000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 2,data));

        mp.add(4000000, address(this), abi.encodeWithSignature("acculm(uint256)", 10));
        mp.add(5000000, address(this), abi.encodeWithSignature("acculm(uint256)", 22));

        require(mp.length() == 5);
        require(container.length() == 0);

        mp.run();      
        require(container.length() == 8);       
        assert(results[0] == keccak256("0"));
        assert(results[1] == keccak256("1"));
    }

    function acculm(uint256 amount) public {
       u256Comulative.add(amount);
       container.push(true);
    }

    function hasher(uint256 idx, bytes memory data) public {
       container.push(true);
       results[idx] = keccak256(data);
       atomic.deferred(300000, address(this), abi.encodeWithSignature("hasherDeferred()"));
    }

    function hasherDeferred() public {
       container.push(true);
       container.push(true);
       container.push(true);
       results[0] = keccak256("0");
       results[1] = keccak256("1");
    }

    function PostCheck() public  {
        assert(results[0] == keccak256("0"));
        assert(results[1] == keccak256("1"));
        require(u256Comulative.get() == 32);
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

contract ConflictInThreadsFixedLengthTest {  
     uint256[2] results;
     function call() public  { 
       results[0] = 100;
       results[1] = 200;
       Threading mp = new Threading(2);
       mp.add(100000, address(this), abi.encodeWithSignature("updater(uint256)", 11));
       mp.add(400000, address(this), abi.encodeWithSignature("updater(uint256)", 33));
       mp.add(400000, address(this), abi.encodeWithSignature("updater(uint256)", 55));
       mp.run();     
       require(results[0] == 155); 
       require(results[1] == 255); 
    }

    function updater(uint256 num) public {
         results[0] += num;
         results[1] += num;
    }
}




