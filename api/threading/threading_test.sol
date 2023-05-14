pragma solidity ^0.5.0;

import "./Threading.sol";

contract ThreadingTest {
    function call() public  { 
       bytes memory data = "0x60298f78cc0b47170ba79c10aa3851d7648bd96f2f8e46a19dbc777c36fb0c00";

       Threading mp = new Threading();
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       mp.add(address(this), abi.encodeWithSignature("hasher(bytes)", data));
       assert(mp.length() == 2);

       mp.del(0);
       assert(mp.length() == 1);

       (bool success,) = address(address(0x90)).call(abi.encodeWithSignature("run()", 1));   
       assert(success);

       (,bytes memory hash) = mp.get(0);
       bytes32 hash32 = bytesToBytes32(hash); 
       assert(hash32 == keccak256(data));

       assert(mp.length() == 1);
       mp.clear();
       assert(mp.length() == 0);
    }

    function callBasic() public  {      
        bytes memory byteArray3 = new bytes(5);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray3[i] = 0x42;
        }

        bytes memory byteArray = new bytes(15);
        for (uint  i = 0; i < byteArray3.length; i ++) {
            byteArray[i] = 0x52;
        }

       bytes memory callArg = abi.encodeWithSignature("hasher(address,bytes)", address(this),byteArray);
       (bool success0,) = address(this).call(callArg);
       assert(success0);

       (success0,) = address(address(0x90)).call(abi.encodeWithSignature("run()", address(this), callArg));   
       assert(success0);
    }

    function localTester (address addr, bytes memory args) public returns(bool, bytes memory)  {
        (bool success, bytes memory id) = address(addr).call(args);
        return (success, id);
    }

    function addJobTester (bytes memory encoded) public returns(bool, bytes memory)  {
        (address addr, bytes memory funcCall) = abi.decode(encoded, (address,bytes));
        (bool success0, bytes memory id) = addr.call(funcCall);        
        assert(success0);
        assert(abi.decode(id, (uint256)) == 112);        
        return (true, funcCall);
    }

    function hasher(bytes memory data) pure public returns(bytes32){
      return keccak256(data);
    }

    function hasher() pure public {}

    function bytesToBytes32(bytes memory b) private pure returns (bytes32) {
        bytes32 out;
        for (uint i = 0; i < 32; i++) {
            out |= bytes32(b[i] & 0xFF) >> (i * 8);
        }
        return out;
    }
}