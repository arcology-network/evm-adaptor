pragma solidity ^0.5.0;

import "./ConcurrentLibInterface.sol";

contract DynamicArrayTest {
    DynamicArray constant darray = DynamicArray(0x84);

    event LogElement(bytes, bool);

    constructor() public {
        darray.create("a-very-long-id-that-is-longer-than-32-bytes", uint256(ConcurrentLib.DataType.BYTES));
        darray.create("short-id", uint256(ConcurrentLib.DataType.BYTES));
    }

    function push(bytes calldata elem) external {
        darray.pushBack("a-very-long-id-that-is-longer-than-32-bytes", elem);
    }

    function tryPop() external returns(bytes memory) {
        bytes memory elem;
        bool ok;
        (elem, ok) = darray.tryPopBackBytes("a-very-long-id-that-is-longer-than-32-bytes");
        emit LogElement(elem, ok);
        return elem;
    }

    function push2(bytes calldata elem) external {
        darray.pushBack("short-id", elem);
    }

    function tryPop2() external returns(bytes memory) {
        bytes memory elem = darray.popFrontBytes("short-id");
        emit LogElement(elem, true);
        return elem;
    }
}