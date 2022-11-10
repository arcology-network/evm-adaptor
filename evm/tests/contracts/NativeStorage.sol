pragma solidity ^0.5.0;

contract NativeStorage {
    struct S {
        uint256 x;
        uint256 y;
    }

    S public s;

    function accessX() public {
        s.x = 1;
    }

    function accessY() public {
        s.y = 2;
    }
}