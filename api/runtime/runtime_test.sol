// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "./Runtime.sol";
import "../threading/Threading.sol";
import "../noncommutative/bool/Bool.sol";
import "../commutative/u256/U256Cumulative.sol";



contract LocalizerTest { 
    uint256 v;
    uint[2] arr;
    mapping(uint256 => uint256) public data;
    Runtime atomic = new Runtime(); 
    constructor() { 
        atomic.localize(3);
        data[2] = 10; 
        data[11] = 21;   
        data[12] = 22;   
        data[152] = 20;   
        data[188] = 20;   

      arr[0] = 41;
      arr[1] = 51;
      v = 61;
      require(arr[0] == 41);
      require(arr[1] == 51);
    }
}


