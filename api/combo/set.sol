// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;

import "../noncommutative/base/Base.sol";

// contract Map_AU{
//   mapping (address => uint256) map;
//   address[] keyArray;

//   constructor () {
//   }


//   function set(address _key, uint256 _value) public {
//     map[_key] = _value;
//     if (!contains(_key)) {
//       keyArray.push(_key);
//     }
//   }

//   /**
//    * Get the value of a key
//    * @param _key The key
//    * @return The value of the key, and a bool indicating if the key exists
//    */
//   function get(address _key) public view returns (uint256, bool) {
//     if (!contains(_key)) {
//       return (0, false);
//     }
//     return (map[_key], true);
//   }

//   /**
//    * Check if the key exists in the map
//    * @param _key The key to check
//    */
//   function contains(address _key) public view returns (bool) {
//     // if the value is not zero, then the key exists
//     if (map[_key] != 0) {return true;}

//     // else, let's check if the key exists in the keyArray array
//     for (uint256 i = 0; i < keyArray.length; i++) {
//       if (keyArray[i] == _key) {return true;}
//     }

//     // if we reach here, then the key doesn't exist
//     return false;
//   }

//   /**
//    * This function removes a key from the map. 
//    * It returns true if the key was removed, and false if the key doesn't exist
//    * @param _key The key to remove. 
//    * @return true if the key was removed, false if the key doesn't exist
//    * 
//    */
//   function remove(address _key) public returns (bool) {
//     if (!contains(_key)) {return false;}
//     // create a new keyArray array, remove key from it, and set it as the new keyArray array
//     address[] memory newkeyArray = new address[](keyArray.length - 1);
//     uint256 j = 0;
//     for (uint256 i = 0; i < keyArray.length; i++) {
//       if (keyArray[i] != _key) {
//         newkeyArray[j] = keyArray[i];
//         j++;
//       }
//     }
//     // set the new keyArray array
//     keyArray = newkeyArray;

//     // remove the key from the map
//     delete map[_key];
//     return true;
//   }

//   /**
//    * Return the size of the map
//    */
//   function size() public view returns (uint256) {
//     return keyArray.length;
//   }

//   function keys() public view returns (address[] memory) {
//     return keyArray;
//   }
// }

contract U256Set is Base {
   mapping (uint256 => uint256) map;
    
   constructor() Base(address(0x84)) {}

   function exist(uint256 key) public virtual returns(bool) { //9e c6 69 25
        return Base.find(abi.encode(key)) < type(uint256).max;
    }

    function insert(uint256 elem) public { // 80 26 32 97
        Base.insert(Base.rand(), abi.encode(elem));  
    }

    function get(uint256 elem) public virtual { //9e c6 69 25
         return abi.decode(Base.getElem(idx), (bool));  
    }    

    function length(uint256 idx) public virtual  returns(bool)  { // 31 fe 88 d0
        return abi.decode(Base.getElem(idx), (bool));  
    }
}
