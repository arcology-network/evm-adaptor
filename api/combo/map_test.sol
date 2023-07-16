// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Map.sol";
import "../multiprocess/Multiprocess.sol";

contract MapTest {
    U256Map map = new U256Map();
    constructor() {     
        require(map.length() == 0); 
        map.set(10, 100);
        map.set(11, 111);
        require(map.length() == 2); 

        require(!map.exist(0));
        require(map.exist(10)); 
        require(map.exist(11)); 

        (,uint256 v) = map.get(11);
        require(v == 111); 

        (,v) = map.get(10);
        require(v == 100); 

        map.del(10);
        require(map.length() == 1); 

        require(map.exist(11));
        map.del(11);
        require(map.length() == 0); 

        require(!map.exist(10));
        require(!map.exist(11));
    }
}

contract ConcurrenctMapTest {
    U256Map map = new U256Map();
    function call() public  { 
       Multiprocess mp = new Multiprocess(2); 
       mp.push(abi.encode(50000, address(this), abi.encodeWithSignature("assigner(uint256)", 11)));
       mp.push(abi.encode(50000, address(this), abi.encodeWithSignature("assigner(uint256)", 33)));
       require(mp.length() == 2);

        (,uint256 v) = map.get(11);
        require(v == 110); 

        (,v) = map.get(33);
        require(v == 330); 

        map.del(11);
        require(map.length() == 1); 
        require(!map.exist(11));

        (,v) = map.get(33);
        require(v == 330); 
        map.del(33);
        require(map.length() == 0); 

        require(!map.exist(11));
        require(!map.exist(33));
    }

    function assigner(uint256 v)  public {
        map.set(v, v * 10);
    }
} 