// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.19;
import "../runtime/Runtime.sol";

contract BaseLinearTest is Runtime{    
    address constant public API = address(0x84); 

    uint[] public arr2 = [1, 2, 3];
    bytes private id;

    // event logMsg(string message);

    constructor() {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new()"));       
        require(success, "Bytes.New() Failed");
        id = data;
 
        bytes memory byteArray = new bytes(75);
        for (uint  i = 0; i < 75; i ++) {
            byteArray[i] = 0x41;
        }
    
        require(peek() == 0);  
        require(length() == 0); 
        push(byteArray);  
        push(byteArray);          
        require(length() == 2); 
        require(peek() == 0);  

        bytes memory stored = getIndex(1);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        bytes memory elems = new bytes(5);
        for (uint  i = 0; i < elems.length; i ++) {
            elems[i] = 0xaa;
        }
        setIndex(1, elems);
       
        stored = getIndex(0);
        require(stored.length == byteArray.length);
        for (uint  i = 0; i < byteArray.length; i ++) {
            require(stored[i] == byteArray[i]);
        }

        stored = getIndex(1);
        require(stored.length == elems.length); 
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }

        // stored = pop();
        stored = getIndex(length() - 1);
        delIndex(length() - 1);
        // stored = pop();
        // require(keccak256(temp) == keccak256(stored));
        for (uint  i = 0; i < elems.length; i ++) {
            require(stored[i] == elems[i]);
        }
        require(length() == 1); 

        // clear();
        require(peek() == 0);  
        require(length() == 1); 
    }

    function call() public{ 
        require(peek() == 1); 
        pop();
        require(peek() == 1); 
    }

    function peek() public returns(uint256) {
        (,bytes memory data) = address(API).call(abi.encodeWithSignature("peek()"));
        if (data.length > 0) {
            return abi.decode(data, (uint256));   
        }
        return 0;    
    }

    function length() public returns(uint256) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));
        require(success, "Bytes.length() Failed");
        return  abi.decode(data, (uint256));
    }

    function pop() public returns(bytes memory) {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("pop()"));
        require(success, "Bytes.pop() Failed");
        return abi.decode(data, (bytes)); 
    }

    function push(bytes memory elem) public {
        setKey(uuid(), (elem));
    }   
    
    function setIndex(uint256 idx, bytes memory encoded) public { // 7a fa 62 38
        address(API).call(abi.encodeWithSignature("setIndex(uint256,bytes)", idx, encoded));     
    }

    function setKey(bytes memory key, bytes memory elem) public {
        address(API).call(abi.encodeWithSignature("setKey(bytes,bytes)", key, elem));
    }

    function delIndex(uint256 idx) public { // 7a fa 62 38
        address(API).call(abi.encodeWithSignature("delIndex(uint256)", idx));     
    }

    function delKey(bytes memory key) public {
        address(API).call(abi.encodeWithSignature("delKey(bytes)", key));
    }


    function getIndex(uint256 idx) public virtual returns(bytes memory) { // 31 fe 88 d0
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getIndex(uint256)", idx));
        if (success) {
            return abi.decode(data, (bytes));  
        }
        return data;
    }

    function getKey(bytes memory key) public returns(bytes memory)  {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getKey(bytes)", key));
        if (success) {
            return abi.decode(data, (bytes));  
        }
        return data;
    }

    // function getIndex(uint256 idx) public returns(bytes memory)  {
    //     (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getIndex(uint256)", idx));
    //     if (success) {
    //         return abi.decode(data, (bytes)); 
    //     }
    //     return "";
    // }

    // function setIndex(uint256 idx, bytes memory elem) public {
    //     (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("setIndex(uint256,bytes)", idx, elem));
    //     require(success, "Bytes.setIndex() Failed");
    // }

    // function getKey(uint256 idx) public returns(bytes memory)  {
    //     (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getKey(uint256)", idx));
    //     require(success, "Bytes.getIndex() Failed");
    //     return abi.decode(data, (bytes));  
    // }

    // function setKey(uint256 idx, bytes memory elem) public {
    //     (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("setKey(uint256,bytes)", idx, elem));
    //     require(success, "Bytes.setIndex() Failed");
    // }

    function clear() public {
        (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("clear()"));
        require(success, "Bytes.setIndex() Failed");
    }
}


// contract BaseNonLinearTest {    
//     address constant public API = address(0x84); 

//     uint[] public arr2 = [1, 2, 3];
//     bytes private id;

//     event logMsg(string message);

//     constructor() {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("new()"));       
//         require(success, "Bytes.New() Failed");
//         id = data;
 
//         bytes memory byteArray = new bytes(75);
//         for (uint  i = 0; i < 75; i ++) {
//             byteArray[i] = 0x41;
//         }
    
//         require(peek() == 0);  
//         require(length() == 0); 
//         push(byteArray);  
//         push(byteArray);          
//         require(length() == 2); 
//         require(peek() == 0);  

//         bytes memory stored = getIndex(1);
//         require(stored.length == byteArray.length);
//         for (uint  i = 0; i < byteArray.length; i ++) {
//             require(stored[i] == byteArray[i]);
//         }

//         bytes memory elems = new bytes(5);
//         for (uint  i = 0; i < elems.length; i ++) {
//             elems[i] = 0xaa;
//         }
//         setIndex(1, elems);
       
//         stored = getIndex(0);
//         require(stored.length == byteArray.length);
//         for (uint  i = 0; i < byteArray.length; i ++) {
//             require(stored[i] == byteArray[i]);
//         }

//         stored = getIndex(1);
//         require(stored.length == elems.length); 
//         for (uint  i = 0; i < elems.length; i ++) {
//             require(stored[i] == elems[i]);
//         }

//         stored = pop();
//         for (uint  i = 0; i < elems.length; i ++) {
//             require(stored[i] == elems[i]);
//         }
//         require(length() == 1); 

//         // clear();
//         require(peek() == 0);  
//         require(length() == 1); 
//     }

//     function call() public{ 
//         require(peek() == 1); 
//         pop();
//         require(peek() == 1); 
//     }

//     function peek() public returns(uint256) {
//         (,bytes memory data) = address(API).call(abi.encodeWithSignature("peek()"));
//         if (data.length > 0) {
//             return abi.decode(data, (uint256));   
//         }
//         return 0;    
//     }

//     function length() public returns(uint256) {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("length()"));
//         require(success, "Bytes.length() Failed");
//         return  abi.decode(data, (uint256));
//     }

//     function pop() public returns(bytes memory) {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("pop()"));
//         require(success, "Bytes.pop() Failed");
//         return abi.decode(data, (bytes)); 
//     }

//     function push(bytes memory elem) public {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("push(bytes)", elem));
//         require(success, "Bytes.push() Failed");
//     }   

//     function getKey(uint256 idx) public returns(bytes memory)  {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("getKey(uint256)", idx));
//         require(success, "Bytes.getIndex() Failed");
//         return abi.decode(data, (bytes));  
//     }

//     function setKey(uint256 idx, bytes memory elem) public {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("setKey(uint256,bytes)", idx, elem));
//         require(success, "Bytes.setIndex() Failed");
//     }

//     function clear() public {
//         (bool success, bytes memory data) = address(API).call(abi.encodeWithSignature("clear()"));
//         require(success, "Bytes.setIndex() Failed");
//     }
// }
