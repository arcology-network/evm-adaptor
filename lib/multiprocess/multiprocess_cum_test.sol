// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Multiprocess.sol";
import "../u256cum/U256Cum.sol";
import "../bool/Bool.sol";

contract MixedRecursiveMultiprocessTest {
    Bool container = new Bool();
    uint256[2] results;
    U256Cumulative cumulative = new U256Cumulative(0, 100);  
    U256Cumulative cumulative2 = new U256Cumulative(50, 80);  

    Multiprocess mp;
    function call() public {
		mp = new Multiprocess(1);
        cumulative.add(50);
        container.push(true);
        mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
        mp.push(abi.encode(9999999, address(this), abi.encodeWithSignature("add()"))); // Only one will go through
        mp.run();
        require(container.length() == 3);

        require(results[0] == 11);
        require(results[1] == 12);
        require(container.length() == 3);
        require(cumulative.get() == 55);
        require(cumulative2.get() == 70);
    } 

    function add() public { //9e c6 69 25
        cumulative.add(10);
        Multiprocess mp2 = new Multiprocess(1); 
        mp2.push(abi.encode(11111111, address(this), abi.encodeWithSignature("add2()")));
        mp2.run();
        container.push(true);              
    }  

    function add2() public {
        cumulative.sub(5);
        cumulative2.add(70);
        results[0] = 11;
        results[1] = 12;
        container.push(true);
    }  
}

contract ParallelCumulativeU256 {
	   U256Cumulative cumulative = new U256Cumulative(0, 100); 
		constructor() {
			require(cumulative.peek() == 0);
			cumulative.add(1);
			cumulative.sub(1);
			require(cumulative.peek() == 0);
		}
	
		function call() public {
			require(cumulative.peek() == 0);
	
			Multiprocess mp = new Multiprocess(1);
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
	
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));   
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1)));
			mp.run();
			require(cumulative.get() == 5);
	
			mp.clear();
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1)));
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
			mp.run();
			require(cumulative.get() == 6);
	
			mp.clear();
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1)));
			mp.run();
			require(cumulative.get() == 5);
	
			mp.clear();
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
			mp.run();
			require(cumulative.get() == 7);      
			require(cumulative.peek() == 0);
	
			mp.clear();
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 50))); // 7 + 50 < 100 => 57
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 50))); // 7 + 50 + 50  > 100 still 57 
			mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 1))); // 7 + 50 + 1  < 100 => 58  
			mp.run();  
	
			require(cumulative.get() == 58);
		}
	
		function call1() public {
		    Multiprocess mp = new Multiprocess(1);
		    mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
		    mp.run();
		    require(cumulative.get() == 2);   
	
		    mp.clear();
		    mp.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 1)));
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
	
	contract ThreadingCumulativeU256SameMpMulti {
		U256Cumulative cumulative = new U256Cumulative(0, 100);     
		function call() public {
			Multiprocess mp1 = new Multiprocess(2);
			mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
			mp1.run();
			mp1.clear();
	 
			mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("add(uint256)", 2)));
			mp1.run(); 
			mp1.clear(); 
	
			mp1.push(abi.encode(200000, address(this), abi.encodeWithSignature("sub(uint256)", 2)));
			mp1.run();   
	
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
	
