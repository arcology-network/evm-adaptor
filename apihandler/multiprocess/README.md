
# Multiprocessor
- [Multiprocessor](#multiprocessor)
  - [1. What is the MP](#1-what-is-the-mp)
    - [1.1. Workflow](#11-workflow)
    - [1.2. Problems with Parallel Deployment](#12-problems-with-parallel-deployment)
  - [2. Solution](#2-solution)
    - [2.1. Caller](#21-caller)
    - [2.2. Applying the Nonce Offset to the Caller](#22-applying-the-nonce-offset-to-the-caller)
    - [2.3. Offsetting Nonce vs Offsetting Increment](#23-offsetting-nonce-vs-offsetting-increment)
    - [2.4. Calculating the Nonce Offset](#24-calculating-the-nonce-offset)
    - [2.5. Factors for Calculating the Nonce Offset](#25-factors-for-calculating-the-nonce-offset)
    - [2.6. Factors to use for Calculating the Nonce Offset](#26-factors-to-use-for-calculating-the-nonce-offset)
      - [2.6.1. MP Address](#261-mp-address)
      - [2.6.2. Why not MP Nonce](#262-why-not-mp-nonce)
  - [3. Clean up the Nonce Offset](#3-clean-up-the-nonce-offset)
  - [4. Example](#4-example)
  - [5. To Do](#5-to-do)

The multiprocessor is a great tool for creating thread-like sub transactions in Solidty. It is allows contracts to process multiple sub transactions in parallel. It is espicallly useful for compuationally intensive tasks. However, the introduction of the multiprocessor has brought forth a few issues that need to be addressed.

## 1. What is the MP

It is logical to think MP == EOA in the sense that it can initiate transactions. However, it is not the case. The MP is still a contract, and it cannot pay the gas fee. The MP is more like a sub transaction creation proxy that can create sub transactions on behalf of the caller. The caller is the one who owns the sub transactions. The caller of the MP should be the one who pays.

### 1.1. Workflow

The MP is a contract that can create sub transactions. The MP is called by the caller. The caller passes the sub transactions to the MP. The MP then creates the sub transactions and sends them to the transaction pool. The sub transactions are then processed in parallel.

### 1.2. Problems with Parallel Deployment

One of the main issues with deploying multiple contracts is the potential for conflicts. This occurs because, when deploying contracts in multiprocessor-created transactions, 
a fresh snapshot is created based on the current snapshot. When multiple sub transactions are created in parallel, they all share the same snapshot. 

If they deploy contracts, they will deploy them to the same address. This is because:

- They all share the same caller address
- The initial nonce of the caller address is the same for all transactions, and as the nonce is ncremented by 1 for each transaction, they all end up with the **same caller nonce**.

The target contract address is calculated as based on caller address + nonce. Because the caller address and nonce are the same for all transactions, the target contract address 
is the same for all transactions. **Consequently, all these transaction are all going to conflict with each other.**

## 2. Solution

Thera are a few solutions to this problem.

1. The simplest solution is banning the deployment of contracts in the MP. This is not a good solution because it limits the functionality of the MP. The MP should be able to deploy contracts. 

2. The second solution is not to care about the conflicts. If the happen, they happen. This, again, is not a good solution because it will limit the functionality of the MP.

3. The third solution is to chang the rules for deploying contracts. Currently, a contract's deployment address is based on the caller's address and nonce. It is possible to decouple the address from the caller's address and nonce. This has significant implications, not only for the MP but also for the entire Ethereum network. It is not a viable solution."

3. The third solution is to disperse the nonce value for each MP created sub transaction. This can be done either by offsetting the initial nonce or by incrementing the nonce by a pseudo-random number. This way, the transactions are deployed to different addresses.

### 2.1. Caller

```solidity
contract Deployer {
    U256 array; 

    constructor() { 
       Multiprocess mp = new Multiprocess(1); 
       mp.push(2500000, address(this), abi.encodeWithSignature("init()"));
       require(mp.length() == 1);
       mp.run();
    }

    function init() public {
        array = new U256();
    }
} 
```

### 2.2. Applying the Nonce Offset to the Caller

This solution is preferred because it does not limit the functionality of the MP. The MP can still deploy contracts. The only thing that changes is the initial nonce. The initial nonce is offset by a pseudo-random number. This number is passed to the MP by the caller. The caller can generate this number using a hash function. The hash function can take the caller address and the current block number as input. This way, the number is unique for each MP-created transaction.

### 2.3. Offsetting Nonce vs Offsetting Increment

We can either offset the initial nonce or the increment by a pseudo-random number. If we offset the initial nonce, the increment remains the same. If we offest the increment, the initial nonce will be the same for all transactions. The only thing that changes is the increment. 

Comparing the two methods, offsetting the initial nonces is the better solution. This is because if we offset the increments, the increments generated by multiple transactions will need to be add up to the initial nonce to get the new nonce. 

- If the offset is too small, the increments are still likely to cause conflicts. 
  
- If the offset is too large, the increments are likely to be too large. When adding to the initial nonce, the new nonce may even overflow.


If we offset the initial nonce, the increments are still going to be comply with current Ethereum rules. The offset nonce can be removed after the sub transactions are processed. The final nonce is still identical to the as they were used with the standard deployment method. 

 
### 2.4. Calculating the Nonce Offset

The nonce offset should be both unique and predictable. It should be unique so that the transactions are not deployed to the same address. It should be predictable so that the caller can predict the address of the deployed contract. 

### 2.5. Factors for Calculating the Nonce Offset

Below is a list of factors that can be used to calculate the nonce offset:

1. The caller's address
2. The caller's nonce
3. Transaction Hash
4. Multiprocessor's address
5. Multiprocessor's nonce
6. API/Thread ID

>>**Depth** isn't suitable for calculating the nonce offset. The depth refers to the call depth of the MP. It is the number of levels of 
>>nested MP calls. The depth is the same for all transactions created by the different MPs at the same level.

![alt](/img/mp.png)

### 2.6. Factors to use for Calculating the Nonce Offset

Not all the available factors are needed to calculate the nonce offset. The caller's address and nonce are the most important factors. As EVM uses the caller's address and nonce to calculate the contract's deployment address. After all, the MPs are deploying contracts on behalf of the caller. 

- MP Address = (Caller Address + Caller Nonce)
- API/Thread ID

#### 2.6.1. MP Address

 The MP's address is important it is key to distinguish a deployement by MP or a standard deployment. In the diagram above, without the MP address, it is impossible to distinguish the deployment `Deploy 0` and `Deploy 2`.

#### 2.6.2. Why not MP Nonce

The MP's nonce isn't an orthogonal factor in this case. 

* Different MP instances will have different deployment addresses. This is guaranteed by the either of 
    * The Caller Address + Caller Nonce in normal deployment
    * The Caller Address + Caller Nonce and the Thread ID in MP deployment, when the former is the identical across different MP created subtransactions.

* Reusing the same MP instance is not a problem because the deployment always cause the caller nonce to increment. In the example below the caller's nonce is incremented by 1 after the first deployment in the first batch. 


```solidity
// First batch
Multiprocessor mp = Multiprocessor(1);
mp.Push(...)
mp.Run()

// Second batch
mp.Push(...)
mp.Run()
```

## 3. Clean up the Nonce Offset

The nonce offset is only used to change the initial nonce and it will be cleaned up after all the sub transactions are processed. The nonce offsets are only used for the deployment of contracts. Once the contracts are deployed, they is no longer needed. But nonce increments will be kept to update the caller's nonce. The final nonce is the nonce of the caller after all the sub transactions are processed.

## 4. Example

Consider the following example: A multiprocessor creates 2 sub transactions. The initial nonce is 10. The offset is 50000 for the first sub transaction and 10000 for the second sub transaction. 

First sub transaction:
* Nonce: 10
* Offset: 50000
* Nonce with offset: 50000 + 10 = 50010
* Increment: 1

Second sub transaction:
* Nonce: 10
* Offset: 10000
* Nonce with offset: 10000 + 10 = 10010
* Increment: 1

In the deployments, 
* The first sub deployment has a nonce of 50010  + 1 = 50011.
* The second sub deployment has a nonce of 10010  + 1 = 10011.

**After removing the offset nonce, while keeping the nonce increment, the final nonce will be 10 + 1 + 1. It is the same as the standard nonce after two deployments.**


## 5. To Do
VM handler should check the caller to make sure it isn't the proxy contract address.

