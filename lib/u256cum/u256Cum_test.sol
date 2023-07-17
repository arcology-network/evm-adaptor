// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./U256Cum.sol";
import "../multiprocess/Multiprocess.sol";

contract CumulativeU256Test {
    U256Cumulative cumulative ;

    constructor() {    
        cumulative = new U256Cumulative(1, 100);  // [1, 100]
        require(cumulative.min() == 1);
        require(cumulative.max() == 100);

        require(cumulative.add(99));
        
        cumulative.sub(99); // This won't succeed, so still 99
        require(cumulative.get() == 99);


        cumulative.add(1);
        require(cumulative.get() == 100);

        cumulative.sub(100); // This won't succeed either, so still 100
        require(cumulative.get() == 100);

        cumulative.sub(99);
        require(cumulative.get() == 1);

        cumulative = new U256Cumulative(0, 100);  // [0, 100]
        require(cumulative.get() == 0);

        require(cumulative.add(99));
        require(cumulative.get() == 99);
        
        require(cumulative.sub(99));
        require(cumulative.get() == 0);

        require(cumulative.min() == 0);
        require(cumulative.max() == 100);
    }
}


// contract ThreadingCumulativeU256 {
// //    U256Cumulative cumulative = new U256Cumulative(0, 100); 
//     Multiprocess mp = new Multiprocess(1);
//     constructor() {
//         // require(cumulative.peek() == 0);
//         // cumulative.add(1);
//         // cumulative.sub(1);
//         // require(cumulative.peek() == 0);
//     }

//     function call() public {
//         // require(cumulative.peek() == 0);

//         // Multiprocess mp = new Multiprocess(1);
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));

//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));   
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1)));
//         // mp.run();
//         // require(cumulative.get() == 5);

//         // mp.clear();
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1)));
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
//         // mp.run();
//         // require(cumulative.get() == 6);

//         // mp.clear();
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1)));
//         // mp.run();
//         // require(cumulative.get() == 5);

//         // mp.clear();
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         // mp.run();
//         // require(cumulative.get() == 7);      
//         // require(cumulative.peek() == 0);

//         // mp.clear();
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 50))); // 7 + 50 < 100 => 57
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 50))); // 7 + 50 + 50  > 100 still 57 
//         // mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1))); // 7 + 50 + 1  < 100 => 58  
//         // mp.run();  

//         // require(cumulative.get() == 58);
//     }

//     // function call1() public {
//     //     Multiprocess mp = new Multiprocess(1);
//     //     mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//     //     mp.run();
//     //     require(cumulative.get() == 2);   

//     //     mp.clear();
//     //     mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1)));
//     //     mp.run();
//     //     require(cumulative.get() == 1);   
//     // }

//     // function call2() public {
//     //     require(cumulative.get() == 1);
//     // }

//     // function add(uint256 elem) public { //9e c6 69 25
//     //     cumulative.add(elem);
//     // }  

//     // function sub(uint256 elem) public { //9e c6 69 25
//     //     cumulative.sub(elem);
//     // }  
// }

// contract ThreadingCumulativeU256SameMpMulti {
//     U256Cumulative cumulative = new U256Cumulative(0, 100);     
//     function call() public {
//         Multiprocess mp1 = new Multiprocess(2);
//         mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         mp1.run();
//         mp1.clear();
 
//         mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         mp1.run(); 
//         mp1.clear(); 

//         mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
//         mp1.run();   

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

// contract ThreadingCumulativeU256DifferentMPMulti {
//     U256Cumulative cumulative = new U256Cumulative(0, 100);     
//     function call() public {
//         // Multiprocess mp1 = new Multiprocess(2);
//         // mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         // mp1.run();
//         // mp1.clear();

//         // Multiprocess mp2 = new Multiprocess(1);
//         // mp2.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
//         // mp2.run(); 
//         // mp2.clear(); 

//         // Multiprocess mp3 = new Multiprocess(1);
//         // mp3.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
//         // mp3.run();   

//         // add(2);
//         // require(cumulative.get() == 4);
//     }

//     function add(uint256 elem) public { //9e c6 69 25
//         cumulative.add(elem);
//     }  

//     function sub(uint256 elem) public { //9e c6 69 25
//         cumulative.sub(elem);
//     }  
// }