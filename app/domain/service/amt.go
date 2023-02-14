package service

import "math/big"

type AMTservice interface {
	// 總供應量
	TotalSupply() *big.Int
	// 擁有者
	Owner() string
	// 代幣名稱
	Name() string
	// 代幣簡稱
	Symbol() string

	// 查訊地址擁有數量
	BalanceOf(addr string) int64

	// 造幣
	Mint(toAddr string, valueData string, nonce, sigHex string) error
	// 授權轉移數量
	Approve(senderAddr, spenderAddr string, valueData string) error
	// 提領授權的代幣
	TransferFrom(senderAddr, fromAddr, toAddr string, valueData, nonce, sigHex string) error
	// 轉移代幣
	Transfer(senderAddr, toAddr string, valueData, nonce, sigHex string) error
}
