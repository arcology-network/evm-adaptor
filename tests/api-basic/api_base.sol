pragma solidity ^0.5.0;


contract Bytes {
    uint256 constant public BYTES = 3;
    address constant public API = address(0x84); 

    uint[] public arr2 = [1, 2, 3];
    bytes private id;

    event LogMsg(string message);

    constructor  (bool isPersistent) public {
        string memory func = "New(uint256, bool) returns(byte32)"; // ac aa 8d 70 
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature(func, BYTES, isPersistent));       
        require(success, "Bytes.New() Failed");
        id = data;
 
        bytes memory byteArray = new bytes(75);
        for (uint  i = 0; i < 75; i ++) {
            byteArray[i] = 0x41;
        }

        Length(); 
        Push(byteArray, arr2);    
    }

    function Length() public returns(uint256) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Length(bytes) returns(uint256)", id));
        require(success, "Bytes.Length() Failed");
        return 123;
    }

    function Delete(uint256 idx) public { // 80 26 32 97
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Delete(bytes, uint256)", id, idx));
        require(success, "Bytes.Pop() Failed");
    }

    function Pop() public {
        Delete(Length() - 1);
    }

    function Push(bytes memory elem, uint[] memory array) public { //9e c6 69 25
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Push(bytes, bytes)",  id, elem));
        require(success, "Bytes.Push() Failed");
    }   

    function Get(uint256 idx) public returns(bytes memory)  { // ef a3 ab 94
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Get(bytes, uint256) returns(bytes)", id, idx));
        require(success, "Bytes.Get() Failed");
        return data;
    }

    function Set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Set(bytes, idx, bytes)", id, idx, elem));
        require(success, "Bytes.Set() Failed");
    }
}
