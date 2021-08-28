pragma solidity ^0.5.0;

import "./ConcurrentLibInterface.sol";

contract BasicFunction {
    ConcurrentHashMap constant hashmap = ConcurrentHashMap(0x81);

    uint256 public sum = 0;
    event LogSum(uint256);

    constructor() public {
        hashmap.create("balance", int32(ConcurrentLib.DataType.ADDRESS), int32(ConcurrentLib.DataType.UINT256));
    }

    function set(address account, uint256 balance) public {
        uint256 origin = hashmap.getUint256("balance", account);
        hashmap.set("balance", account, balance);
        sum -= origin;
        sum += balance;
    }

    function getSum() public {
        emit LogSum(sum);
    }
}