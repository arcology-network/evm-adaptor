pragma solidity ^0.5.0;

import "./U256Cumulative.sol";
import "./Threading.sol";

contract CumulativeU256Test {
    U256Cumulative cumulative ;

    constructor() public {    
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
    
    function call() public {
        Threading mp = new Threading(1);
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));   
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.run();
        require(cumulative.get() == 5);

        mp.clear();
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 1));
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.add(address(this), abi.encodeWithSignature("sub(uint256)", 2));
        mp.run();
        require(cumulative.get() == 6);

        mp.clear();
        mp.add(address(this), abi.encodeWithSignature("sub(uint256)", 1));
        mp.run();
        require(cumulative.get() == 5);

        mp.clear();
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.run();
        require(cumulative.get() == 7);      
    }

    function call1() public {
        Threading mp = new Threading(1);
        mp.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp.run();
        require(cumulative.get() == 2);   

        mp.clear();
        mp.add(address(this), abi.encodeWithSignature("sub(uint256)", 1));
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
    function testCase1() public {
        Threading mp1 = new Threading(1);
        mp1.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp1.run();

        Threading mp2 = new Threading(1);
        mp2.add(address(this), abi.encodeWithSignature("add(uint256)", 2));
        mp2.run();    

        add(2);
        require(cumulative.get() == 6);
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
