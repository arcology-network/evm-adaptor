// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Threading.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";

contract ThreadingParaHasherTest {
    bytes32[2] results;
    function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading(2);
       mp.add(100000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 0,data));
       mp.add(200000, address(this), abi.encodeWithSignature("hasher(uint256,bytes)", 1,data));
       require(mp.length() == 2);
       mp.run();

       (,bytes memory hash) = mp.get(0);
       assert(bytesToBytes32(hash) == keccak256(data));
       assert(bytesToBytes32(hash) == results[0]);
       assert(bytesToBytes32(hash) == results[1]);

       mp.clear();
       assert(mp.length() == 0);       
    }

    function hasher(uint256 idx, bytes memory data)  public returns(bytes32){
      results[idx] = keccak256(data);
      return keccak256(data);
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

contract ThreadingParaContainerManipulationTest {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading(2);
       mp.add(100000, address(this), abi.encodeWithSignature("appender(uint256,bytes)", 0,data));
       mp.add(200000, address(this), abi.encodeWithSignature("appender(uint256,bytes)", 1,data));
       require(mp.length() == 2);
       mp.run();

       (,bytes memory hash) = mp.get(0);
       assert(bytesToBytes32(hash) == keccak256(data));
       assert(bytesToBytes32(hash) == results[0]);
       assert(bytesToBytes32(hash) == results[1]);

       require(container.length() == 2);
       mp.clear();
       assert(mp.length() == 0);       
    }

    function appender(uint256 idx, bytes memory data)  public returns(bytes32){
      container.push(true);
      results[idx] = keccak256(data);
      return keccak256(data);
    }

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}

contract RecursiveThreadingTest {
    uint256[2] results;
    function call() public {
        results[0] = 1;        
        Threading mp = new Threading(1);
        mp.add(900000, address(this), abi.encodeWithSignature("add(uint256)", 11));
        mp.run();
        require(results[1] == 11);
        require(results[0] == 22);
    } 

    function add(uint256 elem) public { //9e c6 69 25
        Threading mp2 = new Threading(1);     
        mp2.add(400000, address(this), abi.encodeWithSignature("add2(uint256)", elem));
        mp2.run();              
    }  

    function add2(uint256 elem) public { //9e c6 69 25
        results[1] = elem; 
        results[0] = elem * 2; 
    }  
}

contract MaxRecursiveThreadingTest {
    Bool container = new Bool();
    uint256 counter = 0;
    function call() public {     
        Threading mp = new Threading(1);
        mp.add(9900000, address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.add(9900000, address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.run(); 
        // require(counter == 7);
        require(container.length() == 3); // 2 + 4 + 8
    } 

    function add(uint256 elem) public { //9e c6 69  
        if (elem >= 6) {
            counter = elem;
            return; 
        }
        container.push(true);
        Threading mp2 = new Threading(2);     
        mp2.add(5900000, address(this), abi.encodeWithSignature("add(uint256)", elem + 2));
        mp2.add(5900000, address(this), abi.encodeWithSignature("add(uint256)", elem + 2));
        mp2.run();              
    }  
}