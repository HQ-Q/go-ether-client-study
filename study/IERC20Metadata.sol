// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


/**
 * @dev 核心ERC20接口（IERC20Metadata的父接口）
 * 包含转账、授权、余额查询等核心方法
 */
interface IERC20 {
    /**
     * @dev 转账：从调用者地址向to地址转账amount个代币
     * @return 转账是否成功
     */
    function transfer(address to, uint256 amount) external returns (bool);

    /**
     * @dev 授权：允许spender从调用者地址支配amount个代币
     * @return 授权是否成功
     */
    function approve(address spender, uint256 amount) external returns (bool);

    /**
     * @dev 代转账：从from地址向to地址转账amount个代币（需先授权）
     * @return 转账是否成功
     */
    function transferFrom(address from, address to, uint256 amount) external returns (bool);

    /**
     * @dev 查询account地址的代币余额
     */
    function balanceOf(address account) external view returns (uint256);

    /**
     * @dev 查询spender已获得from地址的授权额度
     */
    function allowance(address owner, address spender) external view returns (uint256);

    /**
     * @dev 转账事件（代币转账时触发）
     */
    event Transfer(address indexed from, address indexed to, uint256 value);

    /**
     * @dev 授权事件（授权时触发）
     */
    event Approval(address indexed owner, address indexed spender, uint256 value);
}


/**
 * @dev 标准ERC20代币的元数据接口（继承自IERC20）
 * 包含代币名称、符号、小数位等元数据方法
 */
interface IERC20Metadata is IERC20 {
    /**
     * @dev 返回代币名称（如"USDT"）
     */
    function name() external view returns (string memory);

    /**
     * @dev 返回代币符号（如"USDT"）
     */
    function symbol() external view returns (string memory);

    /**
     * @dev 返回代币小数位（通常为18）
     */
    function decimals() external view returns (uint8);
}
