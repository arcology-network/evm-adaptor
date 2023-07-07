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


contract ParaContainerTest {
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

contract ParaMultiTest {
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

contract ParaMultiTestWithClear {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       Parallel mp = new Parallel(2);
       mp.push(abi.encode(1000000, address(this), abi.encodeWithSignature("appender()")));
       mp.run();
       require(container.length() == 1);    

       mp.clear();
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

contract ThreadingCumulativeU256Multi {
    U256Cumulative cumulative = new U256Cumulative(0, 100);     
    function call() public {
        Parallel mp1 = new Parallel(1);
        mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
        mp1.run();

        Parallel mp2 = new Parallel(1);
        mp2.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
        mp2.run();  

        Parallel mp3 = new Parallel(1);
        mp3.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
        mp3.run();   

        add(2);
        require(cumulative.get() == 4);
    }

    function add(uint256 elem) public { //9e c6 69 25
        cumulative.add(elem);
    }  

    function sub(uint256 elem) public { //9e c6 69 25
        cumulative.sub(elem);
    }  
}


// contract ThreadingMpArrayTest {
//     Bool container = new Bool();
//     Threading[2] mps;
//     bytes32[2] results;
//     function call() public  { 
//        mps[0] = new Threading(2);
//        mps[1] = new Threading(2);

//        mps[0].add(1000000, address(this), abi.encodeWithSignature("appender()"));
//        mps[0].run();
//        require(container.length() == 1);    
//        mps[0].clear();

//        mps[0].add(4000000, address(this), abi.encodeWithSignature("appender()"));
//        mps[0].run();
//        require(container.length() == 2);    

//        mps[1].add(4000000, address(this), abi.encodeWithSignature("appender()"));
//        mps[1].add(4000000, address(this), abi.encodeWithSignature("appender()"));
//        mps[1].run();
//        require(container.length() == 4);  
//     }

//     function appender()  public {
//        container.push(true);
//     }
// }

// contract ThreadingMpArraySubprocessTest {
//     Bool container = new Bool();
//     Threading[2] mps;
//     function call() public {     
//         mps[0] = new Threading(1);
//         mps[1] = new Threading(1);  // This will cause a conflict

//         mps[0].add(9000000, address(this), abi.encodeWithSignature("add()"));
//         mps[0].run();
//         require(container.length() == 1);
//     } 

//     function add() public { //9e c6 69 25
//         container.push(true);
//         mps[1].add(4000000, address(this), abi.encodeWithSignature("add2()"));
//         mps[1].run();              
//     }  

//     function add2() public { //9e c6 69 25
//         container.push(true);
//     }  
// }


// contract ThreadingDeploymentAddressTest {
//     Threading[1] mps;
//     function call() public {     
//         Threading mp = new Threading(1);     
//         mp.push(abi.encode(4000000, address(this), abi.encodeWithSignature("deployer()"));  
//         mp.run();    
//     } 

//     function deployer() public { //9e c6 69 25
//         Threading mp = new Threading(1); 
//     }  
 
// }

// contract RecursiveThreadingTest {
//     uint256[2] results;
//     function call() public {
//         results[0] = 1;        
//         Threading mp = new Threading(1);
//         mp.push(abi.encode(900000, address(this), abi.encodeWithSignature("add(uint256)", 11)); // Only one will go through
//         mp.push(abi.encode(900000, address(this), abi.encodeWithSignature("add(uint256)", 41));
//         mp.run();
//         require(results[1] == 11); // Will cause conflicts
//         require(results[0] == 22);
//     } 

//     function add(uint256 elem) public { //9e c6 69 25
//         Threading mp2 = new Threading(1);     
//         mp2.push(abi.encode(400000, address(this), abi.encodeWithSignature("add2(uint256)", elem));
//         mp2.run();              
//     }  

//     function add2(uint256 elem) public { //9e c6 69 25
//         results[1] = elem; 
//         results[0] = elem * 2; 
//     }  
// }



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
