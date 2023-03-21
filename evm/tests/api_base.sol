pragma solidity ^0.5.0;


contract Bytes {
    uint256 constant public BYTES = 3;
    address constant public API = address(0x84); 

    string private id;

    constructor  (bool isPersistent) public {
        string memory func = "New(uint256, bool) returns(byte32)"; // ac aa 8d 70 
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature(func, BYTES, isPersistent));       
        require(success, "Bytes.New() Failed");
        id = abi.decode(data, (string));
    }

    function Length() public returns(uint256) {  // 58 94 13 33
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Length(string) returns(uint256)", id));
        require(success, "Bytes.Length() Failed");
        return abi.decode(data, (uint256));
    }

    function Delete(uint256 idx) public { // 80 26 32 97
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Delete(string, uint256)", id, idx));
        require(success, "Bytes.Pop() Failed");
    }

    function Pop() public {
        Delete(Length() - 1);
    }

    function Push(bytes memory elem) public { // 17 dc 30 f5
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Push(string, bytes)", id, elem));
        require(success, "Bytes.Push() Failed");
    }   

    function Get(uint256 idx) public returns(bytes memory)  { // ef a3 ab 94
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Get(string, uint256) returns(bytes)", id, idx));
        require(success, "Bytes.Get() Failed");
        return data;
    }

    function Set(uint256 idx, bytes memory elem) public { // 7a fa 62 38
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("Set(string, idx, bytes)", id, idx, elem));
        require(success, "Bytes.Set() Failed");
    }
}
