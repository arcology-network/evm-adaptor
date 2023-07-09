// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Parallel.sol";
import "../commutative/u256/U256Cumulative.sol";
import "../noncommutative/bool/Bool.sol";

contract ParaNativeAssignmentTest {
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
    Bool container2 = new Bool();
    function call() public  { 
       container.push(true);

       Parallel mp = new Parallel(2);
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 3);    
       require(container2.length() == 2);   
       container.push(true);
       require(container.length() == 4);    
    //    container.push(true);        
    }

    function appender()  public {
       container.push(true);
       container2.push(true);
    }
}

contract MultiTempParaTest {
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

contract MultiGlobalParaSingleInUse {
    Bool container = new Bool();

    Parallel mp2 = new Parallel(2);
    Parallel mp = new Parallel(2);
    function call() public  {  
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.run();
       require(container.length() == 2);    
    }

    function appender()  public {
       container.push(true);
       container.push(true);
    }
}

contract MultiGlobalPara {
    Bool container = new Bool();

    Parallel mp = new Parallel(2);
    Parallel mp2 = new Parallel(2);
    function call() public  {  
       mp.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 2);     
      
       mp2 = new Parallel(2);
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.push(abi.encode(4000000, address(this), abi.encodeWithSignature("appender()")));
       mp2.run();
       require(container.length() == 4);  

       container.push(true);
       require(container.length() == 5);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract MultiLocalParaTestWithClear {
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

contract MixedRecursiveParallelizerTest {
    Bool container = new Bool();
    uint256[2] results;
    U256Cumulative cumulative = new U256Cumulative(0, 100);  

    Parallel mp = new Parallel(1);
    function call() public {
        cumulative.add(50);
        container.push(true);
        mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
        mp.run();

        require(results[0] == 11);
        require(results[1] == 12);
        require(container.length() == 3);
        require(cumulative.get() == 55);
    } 

    function add() public { //9e c6 69 25
        cumulative.add(10);
        Parallel mp2 = new Parallel(1); 
        mp2.push(abi.encode(11111111, address(this), abi.encodeWithSignature("add2()")));
        mp2.run();
        container.push(true);              
    }  

    function add2() public {
        cumulative.sub(5);
        results[0] = 11;
        results[1] = 12;
        container.push(true);
    }  
}


// contract MixedRecursiveParallelizerTest {
//     Bool container = new Bool();
//     uint256[2] results;
//     U256Cumulative cumulative = new U256Cumulative(0, 100);  

//     Parallel mp;
//     function call() public {
//         cumulative.add(50);
//         container.push(true);       
//         mp = new Parallel(1);
//         mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
//         mp.run();

//         require(results[0] == 11);
//         require(results[1] == 12);
//         require(container.length() == 3);
//         require(cumulative.get() == 55);
//     } 

//     function add() public { //9e c6 69 25
//         cumulative.add(10);
//         Parallel mp2 = new Parallel(1); 
//         mp2.push(abi.encode(11111111, address(this), abi.encodeWithSignature("add2()")));
//         mp2.run();
//         container.push(true);              
//     }  

//     function add2() public {
//         cumulative.sub(5);
//         results[0] = 11;
//         results[1] = 12;
//         container.push(true);
//     }  
// }
