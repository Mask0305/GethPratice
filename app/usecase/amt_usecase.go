package usecase

import "scp/app/domain/service"

type AMT interface {
	TotalSupply() int64
	// 合約擁有者
	Owner() string
	// 代幣名稱
	Name() string
	// 代幣簡稱
	Symbol() string
	// 授權轉移數量
	Approve(
		senderAddr, spenderAddr string,
		valueData string,
	) error
	// 查詢地址擁有數量
	BalanceOf(addr string) int64
	// 造幣
	Mint(toAddr string, valueData string, nonce, sigHex string) error
	// 提領授權的代幣
	TransferFrom(
		senderAddr, fromAddr, toAddr string,
		valueData string,
		nonce, sigHex string,
	) error
	// 轉移代幣
	Transfer(
		senderAddr, toAddr string,
		valueData string,
		nonce, sigHex string,
	) error
}

type amtUsecase struct {
	amtSrc service.AMTservice
}

func New(amtSrc service.AMTservice) AMT {
	return &amtUsecase{
		amtSrc: amtSrc,
	}
}

// 總供應量
func (a *amtUsecase) TotalSupply() int64 {
	return a.amtSrc.TotalSupply().Int64()
}

// 合約擁有者
func (a *amtUsecase) Owner() string {
	return a.amtSrc.Owner()
}

// 代幣名稱
func (a *amtUsecase) Name() string {
	return a.amtSrc.Name()
}

// 代幣簡稱
func (a *amtUsecase) Symbol() string {
	return a.amtSrc.Symbol()
}

func (a *amtUsecase) BalanceOf(addr string) int64 {
	return a.amtSrc.BalanceOf(addr)
}

// 造幣
func (a *amtUsecase) Mint(toAddr string, valueData string, nonce, sigHex string) error {
	return a.amtSrc.Mint(toAddr, valueData, nonce, sigHex)
}

// 授權轉移數量
func (a *amtUsecase) Approve(
	senderAddr, spenderAddr string,
	valueData string,
) error {
	return a.amtSrc.Approve(senderAddr, spenderAddr, valueData)
}

// 提領授權的代幣
func (a *amtUsecase) TransferFrom(
	senderAddr, fromAddr, toAddr string,
	valueData string,
	nonce, sigHex string,
) error {
	return a.amtSrc.TransferFrom(senderAddr, fromAddr, toAddr, valueData, nonce, sigHex)
}

// 轉移代幣
func (a *amtUsecase) Transfer(
	senderAddr, toAddr string,
	valueData string,
	nonce, sigHex string,
) error {
	return a.amtSrc.Transfer(senderAddr, toAddr, valueData, nonce, sigHex)
}
