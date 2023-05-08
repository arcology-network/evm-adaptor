pragma solidity ^0.5.0;

import "./Multiprocess.sol";


contract String {

    function reverse(bytes memory message) internal pure returns(bytes memory){
        bytes memory reversed = new bytes(message.length);
        for(uint i=0;i<message.length;i++){
            reversed[ message.length - i - 1] = message[i];
        }
        return reversed;
    }

    function testReverseString10k() public {
        string memory message = "0x98765299729916702639311866201502563458355484178623343495093432660206149650039545677028594864142235478903100571966703942682266471199453872255049900266750127402724904010526205528730342205570736334250845";

        Multiprocess mp = new Multiprocess();
        for (int i = 0; i < 10000; i ++) { 
            mp.addJob(address(this), abi.encodeWithSignature("reverse(bytes)", address(this), message));
        }
        mp.run();
    }
}