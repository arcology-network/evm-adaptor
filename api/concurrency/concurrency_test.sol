// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Concurrency.sol";
import "../threading/Threading.sol";


contract ConcurrencyDeferredInThreadingTest {
     function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading(2);
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       require(mp.length() == 2);
       mp.run();

       (,bytes memory hash) = mp.get(0);
       assert(bytesToBytes32(hash) == keccak256(data));

       mp.clear();
       assert(mp.length() == 0);       
    
    
      // this should fail because the deferred call hasn't been processed yet. 
    
    }

    function hasher(bytes memory data)  public returns(bytes32){
       Concurrency concurrency = new Concurrency(); 
       concurrency.uuid(); 

       (bool success, bytes memory native256) = concurrency.deferred(address(this), abi.encodeWithSignature("example()"));
       require(success);     
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }

    function example() pure public returns(bytes32){
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
       return keccak256(data);
    }
}

contract ConcurrencyDeferredTest {
     function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading(2);
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       require(mp.length() == 2);
       mp.run();

       (,bytes memory hash) = mp.get(0);
       assert(bytesToBytes32(hash) == keccak256(data));

       mp.clear();
       assert(mp.length() == 0);       
    }

    function hasher(bytes memory data)  public returns(bytes32){
       Concurrency concurrency = new Concurrency(); 
       concurrency.uuid(); 

       (bool success, bytes memory native256) = concurrency.deferred(address(this), abi.encodeWithSignature("example()"));
       require(success && keccak256(data) == bytesToBytes32(native256));
       return keccak256(data);
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }

    function example() pure public returns(bytes32){
      bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";
      return keccak256(data);
    }
}


