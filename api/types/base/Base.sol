pragma solidity ^0.5.0;

contract Base {
    address constant public API = address(0x84);    
    bytes private ctrn;

    event logMsg(string message);

    constructor () public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new()"));       
        require(success);
        ctrn = data; 
    }

    function id() public view returns(bytes memory) {
        return ctrn;
    }

    function length() public returns(uint256) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length(bytes)", ctrn));
        require(success);
        return abi.decode(data, (uint256));
    }

    function pop() public returns(bytes memory) { // 80 26 32 97
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("pop()", ctrn));
        require(success);
        return data; 
    }

    function push(bytes memory encoded) public { //9e c6 69 25
        (bool success, bytes memory data) = address(API).call(encoded);
        require(success);
    }   

    function get(uint256 idx) public returns(bytes memory)  { // 31 fe 88 d0
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("get(bytes,uint256)", ctrn, idx));
        require(success);
        return data;  
    }

    function set(bytes memory encoded) public { // 7a fa 62 38
        (bool success, bytes memory data) = address(API).call(encoded);
        require(success);
    }
    
    function log(bytes memory elem) public { // 7a fa 62 38
        address(API).call(abi.encodeWithSignature("log(bytes)", id(), elem));     
    }
}
