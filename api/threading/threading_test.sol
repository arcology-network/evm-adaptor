// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Threading.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";

contract ThreadingNativeArrayTest {
    uint256[2] results;
    function call() public  { 
    //    bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading(2);
       mp.add(50000, address(this), abi.encodeWithSignature("assigner()"));
       mp.add(50000, address(this), abi.encodeWithSignature("assigner()"));
       require(mp.length() == 2);
       mp.run();

       assert(results[0] == 11);
    //    assert(results[1] == 11);
   }

    function assigner()  public {
        results[0] = 1 + 10;
    }
}


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

contract ThreadingFixedLengthWithConflictTest {  
     uint256[2] results;
     function call() public  { 
       results[0] = 100;
       results[1] = 200;
       Threading mp = new Threading(2);
       mp.add(400000, address(this), abi.encodeWithSignature("updater(uint256)", 11));
       mp.add(400000, address(this), abi.encodeWithSignature("updater(uint256)", 33));
       mp.add(400000, address(this), abi.encodeWithSignature("updater(uint256)", 55));
       mp.run();     
       require(results[0] == 111);  // 11 and 33 will be reverted due to conflicts
       require(results[1] == 211); 
    }

    function updater(uint256 num) public {
         results[0] += num;
         results[1] += num;
    }
}

contract ThreadingParaContainerManipulationTest {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       container.push(true);
       Threading mp = new Threading(2);
       mp.add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mp.add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mp.run();

       require(container.length() == 5);    
    }

    function appender()  public {
       container.push(true);
       container.push(true);
    }
}

// contract ThreadingCumulativeU256Multi {
//     U256Cumulative cumulative = new U256Cumulative(0, 100);     
//     function call() public {
//         Threading mp1 = new Threading(1);
//         mp1.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
//         mp1.run();

//         Threading mp2 = new Threading(1);
//         mp2.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
//         mp2.run();  

//         Threading mp3 = new Threading(1);
//         mp3.add(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2));
//         mp3.run();   

//         add(2);
//         require(cumulative.get() == 4);
//     }

//     function add(uint256 elem) public { //9e c6 69 25
//         cumulative.add(elem);
//     }  

//     function sub(uint256 elem) public { //9e c6 69 25
//         cumulative.sub(elem);
//     }  
// }

contract ThreadingMultiMPsTest {
    Bool container = new Bool();
    bytes32[2] results;
    function call() public  { 
       Threading mp = new Threading(2);
       mp.add(1000000, address(this), abi.encodeWithSignature("appender()"));
       mp.run();
       require(container.length() == 1);    

       mp.clear();
       mp.add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mp.run();
       require(container.length() == 2);    

       Threading mp2 = new Threading(2);
       mp2.add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mp2.add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mp2.run();
       require(container.length() == 4);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract ThreadingMpArrayTest {
    Bool container = new Bool();
    Threading[2] mps;
    bytes32[2] results;
    function call() public  { 
       mps[0] = new Threading(2);
       mps[1] = new Threading(2);

       mps[0].add(1000000, address(this), abi.encodeWithSignature("appender()"));
       mps[0].run();
       require(container.length() == 1);    
       mps[0].clear();

       mps[0].add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mps[0].run();
       require(container.length() == 2);    

       mps[1].add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mps[1].add(4000000, address(this), abi.encodeWithSignature("appender()"));
       mps[1].run();
       require(container.length() == 4);  
    }

    function appender()  public {
       container.push(true);
    }
}

contract ThreadingMpArraySubprocessTest {
    Bool container = new Bool();
    Threading[2] mps;
    function call() public {     
        mps[0] = new Threading(1);
        mps[1] = new Threading(1);  // This will cause a conflict

        mps[0].add(9000000, address(this), abi.encodeWithSignature("add()"));
        mps[0].run();
        require(container.length() == 1);
    } 

    function add() public { //9e c6 69 25
        container.push(true);
        mps[1].add(4000000, address(this), abi.encodeWithSignature("add2()"));
        mps[1].run();              
    }  

    function add2() public { //9e c6 69 25
        container.push(true);
    }  
}


contract ThreadingDeploymentAddressTest {
    Threading[1] mps;
    function call() public {     
        Threading mp = new Threading(1);     
        mp.add(4000000, address(this), abi.encodeWithSignature("deployer()"));  
        mp.run();    
    } 

    function deployer() public { //9e c6 69 25
        Threading mp = new Threading(1); 
    }  
 
}

contract RecursiveThreadingTest {
    uint256[2] results;
    function call() public {
        results[0] = 1;        
        Threading mp = new Threading(1);
        mp.add(900000, address(this), abi.encodeWithSignature("add(uint256)", 11)); // Only one will go through
        mp.add(900000, address(this), abi.encodeWithSignature("add(uint256)", 41));
        mp.run();
        require(results[1] == 11); // Will cause conflicts
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



// contract MaxRecursiveThreadingTest {
//     uint256[2] results;
//     uint256 counter = 0;
//     function call() public {     
//         Threading mp = new Threading(2);
//         mp.add(9900000, address(this), abi.encodeWithSignature("add(uint256)"));
//         mp.add(9900000, address(this), abi.encodeWithSignature("add()"));
//         mp.run(); 
//         require(container.length() ==2); // 2 + 4 + 8
//     } 

//     function add() public { //9e c6 69  
//         container.push(true);
//         Threading mp2 = new Threading(2);     
//         // mp2.add(4400000, address(this), abi.encodeWithSignature("add2()"));
//         // // mp2.add(4400000, address(this), abi.encodeWithSignature("add2()"));
//         // mp2.run();              
//     }  

//     function add2() public { //9e c6 69  
//         container.push(true); 
//     }   
// }
