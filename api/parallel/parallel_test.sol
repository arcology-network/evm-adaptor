// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Parallel.sol";
import "../commutative/u256/U256Cumulative.sol";
import "../noncommutative/bool/Bool.sol";

contract ParaHasherTest {
    uint256[2] results;
    function call() public  { 
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(50000, address(this), abi.encodeWithSignature("assigner(uint256)", 0)));
       mp.push(abi.encode(50000, address(this), abi.encodeWithSignature("assigner(uint256)", 1)));
       require(mp.length() == 2);
       mp.run();

       assert(results[0] == 10);
       assert(results[1] == 11);
    }

    function assigner(uint256 v)  public {
        results[v] = v + 10;
    }
}

contract ParaFixedLengthWithConflictTest {  
     uint256[2] results;
     function call() public  { 
       results[0] = 100;
       results[1] = 200;
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(400000, address(this), abi.encodeWithSignature("updater(uint256)", 11)));
       mp.push(abi.encode(400000, address(this), abi.encodeWithSignature("updater(uint256)", 33)));
       mp.push(abi.encode(400000, address(this), abi.encodeWithSignature("updater(uint256)", 55)));
       mp.run();     
       require(results[0] == 111);  // 11 and 33 will be reverted due to conflicts
       require(results[1] == 211); 
    }

    function updater(uint256 num) public {
         results[0] += num;
         results[1] += num;
    }
}

contract ParaContainerConcurrentPushTest {
    Bool container = new Bool();
    function call() public  { 
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 2);    
    }

    function appender()  public {
       container.push(true);
    }
}

contract MultiParaTest {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 2);     

       Parallel mp2 = new Parallel(2);
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.run();
       require(container.length() == 4);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract MultiParaTestWithClear {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 1);    

       mp.clear();       
       require(mp.length() == 0);   

       mp.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 2);    

       Parallel mp2 = new Parallel(2);
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.run();
       require(container.length() == 4);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract ParallelizerArrayTest {
    Bool container = new Bool();
    Parallel[2] parallelizers;

    function call() public  { 
       parallelizers[0] = new Parallel(2);
       parallelizers[0] .push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       parallelizers[0] .push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       parallelizers[0] .run();
       require(container.length() == 2);  

       parallelizers[1] = new Parallel(2);
       parallelizers[1] .push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       parallelizers[1] .push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       parallelizers[1] .run();
       require(container.length() == 4);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract MultiParaCumulativeU256 {
    U256Cumulative cumulative = new U256Cumulative(0, 100);     
    function call() public {
        Parallel mp1 = new Parallel(1);
        mp1.push(abi.encode(400000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
        mp1.run();

        Parallel mp2 = new Parallel(1);
        mp2.push(abi.encode(400000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
        mp2.run();  

        Parallel mp3 = new Parallel(1);
        mp3.push(abi.encode(400000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
        mp3.run();   

        add(3);
        require(cumulative.get() == 5);
    }

    function add(uint256 elem) public { //9e c6 69 25
        cumulative.add(elem);
    }  

    function sub(uint256 elem) public { //9e c6 69 25
        cumulative.sub(elem);
    }  
}

contract RecursiveParallelizerOnNativeArrayTest {
    uint256[2] results;
    function call() public {
        Parallel mp = new Parallel(1);
        mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
        mp.run();

        require(results[0] == 11);
        require(results[1] == 12);
    } 

    function add() public { //9e c6 69 25
        Parallel mp2 = new Parallel(1); 
        mp2.push(abi.encode(11111111, address(this), abi.encodeWithSignature("add2()")));
        mp2.run();              
    }  

    function add2() public { 
        results[0] = 11;
        results[1] = 12;
    }  
}

contract RecursiveParallelizerOnContainerTest {
    uint256[2] results;
    Bool container = new Bool();
    U256Cumulative cumulative = new U256Cumulative(0, 100);  

    function call() public {
        Parallel mp = new Parallel(1);
        mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
        mp.run();

        require(results[0] == 11);
        require(results[1] == 12);
        require(container.length() == 2);
        require(cumulative.get() == 5);
    } 

    function add() public { //9e c6 69 25
        container.push(true);
        cumulative.add(10);
        Parallel mp2 = new Parallel(1); 
        mp2.push(abi.encode(11111111, address(this), abi.encodeWithSignature("add2()")));
        mp2.run();              
    }  

    function add2() public {
        container.push(true); 
        cumulative.sub(5);
        results[0] = 11;
        results[1] = 12;
    }  
}



// contract MaxRecursiveThreadingTest {
//     uint256[2] results;
//     uint256 counter = 0;
//     function call() public {     
//         Threading mp = new Threading(2);
//         mp.push(abi.encode(9900000, address(this), abi.encodeWithSignature("add(uint256)"));
//         mp.push(abi.encode(9900000, address(this), abi.encodeWithSignature("add()"));
//         mp.run(); 
//         require(container.length() ==2); // 2 + 4 + 8
//     } 

//     function add() public { //9e c6 69  
//         container.push(true);
//         Threading mp2 = new Threading(2);     
//         // mp2.push(abi.encode(4400000, address(this), abi.encodeWithSignature("add2()"));
//         // // mp2.push(abi.encode(4400000, address(this), abi.encodeWithSignature("add2()"));
//         // mp2.run();              
//     }  

//     function add2() public { //9e c6 69  
//         container.push(true); 
//     }   
// }
