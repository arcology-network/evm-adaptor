// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Atomic.sol";
import "../threading/Threading.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";

 // this should fail because the deferred call hasn't been processed yet. 
contract AtomicDeferredInThreadingTest {  
     bytes32[2] results;
     Atomic atomic = new Atomic(); 
     function call() public  { 
        bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
  
        Threading mp = new Threading(2);
        mp.add(800000, address(this), abi.encodeWithSignature("worker(uint256,bytes)", 0,data));
        mp.add(800000, address(this), abi.encodeWithSignature("worker(uint256,bytes)", 1,data));
        require(mp.length() == 2);
        mp.run();
       
        // if (atomic.singleton()) {
            results[0] = keccak256("0");
            results[1] = keccak256("1");
        // }

        mp.clear();
        assert(mp.length() == 0);      
    }

    function worker(uint256 idx, bytes memory data) public {
       results[idx] = keccak256(data);
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
       require(results[0] == 155);  // 11 and 33 will be reverted due to conflicts
       require(results[1] == 255); 
    }

    function updater(uint256 num) public {
         results[0] += num;
         results[1] += num;
    }
}




