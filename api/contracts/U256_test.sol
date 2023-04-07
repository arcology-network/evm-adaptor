pragma solidity ^0.5.0;

import "./U256.sol";

contract U256Test {
    U256 container = new U256();
    U256[] array;

    constructor() public {     
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
        U256 container0 = new U256();
        U256 container1 = new U256();
        U256 container2 = new U256();

        array.push(container0);
        array.push(container1);
        array.push(container2);

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