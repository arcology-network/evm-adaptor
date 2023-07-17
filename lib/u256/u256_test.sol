// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./U256.sol";
import "../multiprocess/Multiprocess.sol";

contract U256DynamicTest {
    U256 container = new U256();
    U256[] array;

    constructor() {     
        require(container.length() == 0); 
    
        container.push(uint256(10));
        container.push(uint256(20));
        container.push(uint256(30));
        container.push(uint256(40));
        require(container.length() == 4); 

        require(container.get(0) == uint256(10));
        require(container.get(1) == uint256(20));
        require(container.get(2) == uint256(30));
        require(container.get(3) == uint256(40));    

        container.set(0, uint256(11));
        container.set(1, uint256(12));
        container.set(2, uint256(13));
        container.set(3, uint256(14));

        require(container.get(0) == uint256(11));
        require(container.get(1) == uint256(12));
        require(container.get(2) == uint256(13));
        require(container.get(3) == uint256(14));

        require(container.pop() == uint256(14));
        require(container.pop() == uint256(13));
        require(container.pop() == uint256(12));
        require(container.pop() == uint256(11));
        require(container.length() == 0); 
        
        // Nested array
        // U256 container0 = new U256();
        // U256 container1 = new U256();
        // U256 container2 = new U256();

        array = new U256[](3);
        array[0] = new U256();
        array[1] = new U256();
        array[2] = new U256();

        // array.push(container0);
        // array.push(container1);
        // array.push(container2);

        U256 _0 = array[0];
        _0.push(101);
        _0.push(102);
        _0.push(103);
        _0.push(104);
        require(_0.length() == 4); 

        require(_0.get(0) == uint256(101));
        require(_0.get(1) == uint256(102));
        require(_0.get(2) == uint256(103));
        require(_0.get(3) == uint256(104));   


        require(_0.pop() == uint256(104));
        require(_0.pop() == uint256(103));
        require(_0.pop() == uint256(102));
        require(_0.pop() == uint256(101));
        require(_0.length() == 0); 

        U256 _1= array[1];
        _1.push(111);
        _1.push(112);
        _1.push(113);
        _1.push(114);

        require(_1.get(0) == uint256(111));
        require(_1.get(1) == uint256(112));
        require(_1.get(2) == uint256(113));
        require(_1.get(3) == uint256(114));   

        require(_1.pop() == uint256(114));
        require(_1.pop() == uint256(113));
        require(_1.pop() == uint256(112));
        require(_1.pop() == uint256(111));
        require(_1.length() == 0); 
    }
}

contract U256ParallelTest {
    U256 container = new U256();

    function call() public  { 
        require(container.length() == 0); 
    
        container.push(uint256(10));
        container.push(uint256(20));
        container.push(uint256(30));
        require(container.length() == 3);

        Multiprocess mp = new Multiprocess(1);
        mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("push(uint256)", 41)));
        mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("push(uint256)", 51)));
        require(mp.length() == 2);
        require(container.length() == 3);

        mp.run();
        require(container.length() == 5);

        require(container.get(0) == uint256(10));
        require(container.get(1) == uint256(20));
        require(container.get(2) == uint256(30));
        require(container.get(3) == uint256(41));   
        require(container.get(4) == uint256(51));  
 
        require(container.pop() == uint256(51));  
        require(container.length() == 4);


        mp.clear();
        mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("get(uint256)", 0)));
        mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("get(uint256)", 1)));
        require(mp.length() == 2);
        mp.run();

        // (bytes memory data) = mp.get(1);
        // require(abi.decode(data, (uint256)) == 20);

        // (data) = mp.get(0);
        // require(abi.decode(data, (uint256)) == 10);

        pop();
        require(container.length() == 3);

        // // Here should be one conflict
        // mp.clear();
        mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("pop()")));
        mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("pop()")));
        mp.run();

        require(container.length() == 1);  // Only one transaction went through, so only one pop() took effect
    }

    function push(uint256 v) public{
        container.push(v);
    }

    function get(uint256 idx) public returns(uint256){
        return container.get(idx);  
    }

    function set(uint256 idx, uint256 v) public {
        return container.set(idx, v);  
    }

    function pop() public {
        container.pop();   
    }
}

contract ArrayParallelTest {
    U256[] array; 

    constructor() {
        array = new U256[](2);
        array[0] = new U256();
        array[1] = new U256();
    }

    function call() public  {
        Multiprocess mp = new Multiprocess(1);
        push(0, 11);
        push(0, 12);

         mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("push(uint256,uint256)", 0, 13)));
         mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("push(uint256,uint256)", 0, 14)));
         mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("push(uint256,uint256)", 1, 51)));
         mp.push(abi.encode(100000, address(this), abi.encodeWithSignature("push(uint256,uint256)", 1, 52)));
        require(mp.length() == 4);
        mp.run();

        require(array[0].length() == 4);
        require(array[1].length() == 2);

        require(array[0].get(0) == 11);
        require(array[0].get(1) == 12);
        require(array[0].get(2) == 13);
        require(array[0].get(3) == 14);

        require(array[1].get(0) == 51);
        require(array[1].get(1) == 52);
    }

    function push(uint256 id, uint256 v) public{
        array[id].push(v);
    }

    function get(uint256 id, uint256 idx) public returns(uint256){
        return array[id].get(idx);  
    }

    function set(uint256 id, uint256 idx, uint256 v) public {
        return array[id].set(idx, v);  
    }

    function pop(uint256 id) public {
        array[id].pop();  
    }
}