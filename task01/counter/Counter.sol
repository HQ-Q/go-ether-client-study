// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

//计数器合约
contract Counter {
    uint256 public count;

    event CountChanged(uint256 newCount);

    constructor() {
        count = 0;
    }
    //增加计数
    function increment() public {
        count += 1;
        emit CountChanged(count);
    }

    //减少计数
    function decrement() public {
        count -= 1;
        emit CountChanged(count);
    }

    //获取当前计数
    function getCount() public view returns (uint256) {
        return count;
    }
}
