// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./U256Cumulative.sol";
import "../../threading/Threading.sol";

contract CumulativeU256Test {
    U256Cumulative cumulative ;

    constructor() {    
        cumulative = new U256Cumulative(1, 100);  // [1, 100]
        require(cumulative.min() == 1);
        require(cumulative.max() == 100);

        require(cumulative.add(99));
        
        cumulative.sub(99);
        require(cumulative.get() == 99);


        cumulative.add(1);
        require(cumulative.get() == 100);

        cumulative.sub(100);
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

contract ThreadingCumulativeU256 {
    U256Cumulative cumulative = new U256Cumulative(0, 100); 

    constructor() {
        require(cumulative.peek() == 0);
        cumulative.add(1);
        cumulative.sub(1);
         require(cumulative.peek() == 0);
    }

    function call() public {
        require(cumulative.peek() == 0);

        Threading mp = new Threading(1);
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));   
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.run();
        require(cumulative.get() == 5);

        mp.clear();
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.add(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2));
        mp.run();
        require(cumulative.get() == 6);

        mp.clear();
        mp.add(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1));
        mp.run();
        require(cumulative.get() == 5);

        mp.clear();
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.run();
        require(cumulative.get() == 7);      
        require(cumulative.peek() == 0);

        mp.clear();
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 50)); // 7 + 50 < 100 => 57
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 50)); // 7 + 50 + 50  > 100 still 57 
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 1)); // 7 + 50 + 1  < 100 => 58  
        mp.run();  

        require(cumulative.get() == 58);
    }

    function call1() public {
        Threading mp = new Threading(1);
        mp.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.run();
        require(cumulative.get() == 2);   

        mp.clear();
        mp.add(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1));
        mp.run();
        require(cumulative.get() == 1);   
    }

    function call2() public {
        require(cumulative.get() == 1);
    }

    function add(uint256 elem) public { //9e c6 69 25
        cumulative.add(elem);
    }  

    function sub(uint256 elem) public { //9e c6 69 25
        cumulative.sub(elem);
    }  
}

contract ThreadingCumulativeU256Multi {
    U256Cumulative cumulative = new U256Cumulative(0, 100);     
    function call() public {
        Threading mp1 = new Threading(1);
        mp1.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp1.run();

        Threading mp2 = new Threading(1);
        mp2.add(200000, address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp2.run();  

        Threading mp3 = new Threading(1);
        mp3.add(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2));
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

contract ThreadingCumulativeU256Recursive {
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
