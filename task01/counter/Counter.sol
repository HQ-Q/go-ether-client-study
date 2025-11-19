// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

//计数器合约
contract Counter {
    uint256 public count;
    constructor() {
        count = 0;
    }
    //增加计数
    function increment() public {
        count += 1;
    }

    //减少计数
    function decrement() public {
        count -= 1;
    }

    //获取当前计数
    function getCount() public view returns (uint256) {
        return count;
    }
}
