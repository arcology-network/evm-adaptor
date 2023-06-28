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

contract RecursiveThreadingTest  {
    U256Cumulative cumulative = new U256Cumulative(0, 100);     
    function testCase1() public {
        Threading mp = new Threading(1);
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.run();
    }

    function add(uint256 elem) public { //9e c6 69 25
        // if (elem > 10) {
        //     cumulative.add(elem);
        //     return;
        // }

        Threading mp = new Threading(1);
        mp.add(200000, address(this), abi.encodeWithSignature("add2(uint256)", 2));
        mp.add(200000, address(this), abi.encodeWithSignature("add2(uint256)", 2));
        mp.run();       
    }  

    function add2(uint256 elem) public { //9e c6 69 25     
    }  
}
