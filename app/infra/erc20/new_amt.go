package erc20

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"scp/app/domain/service"
	"scp/app/infra/erc20/amt"
	"scp/app/infra/ethereum"
	"strings"

	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type amtSrc struct {
	client *ethereum.EthService
	token  *amt.Token
}

func NewAMTtoken(client *ethereum.EthService) service.AMTservice {
	// 合約地址
	tokenAddress := common.HexToAddress("0xA65f533dDcc437F942b0e360C12c0617B2732dFF")

	instance, err := amt.NewToken(tokenAddress, client.Client)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &amtSrc{
		client: client,
		token:  instance,
	}
}

// 總供應量
func (a *amtSrc) TotalSupply() *big.Int {

	n, err := a.token.TotalSupply(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return n
}

// 擁有者
func (a *amtSrc) Owner() string {

	address, err := a.token.Owner(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return address.Hex()
}

// 代幣名稱
func (a *amtSrc) Name() string {

	name, err := a.token.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
		return ""
	}

	return name
}

// 代幣簡稱
func (a *amtSrc) Symbol() string {

	symbol, err := a.token.Symbol(&bind.CallOpts{})

	if err != nil {
		log.Fatal(err)
		return ""
	}

	return symbol
}

// 查訊地址擁有數量
func (a *amtSrc) BalanceOf(addr string) int64 {
	address := common.HexToAddress(addr)

	t, err := a.token.BalanceOf(&bind.CallOpts{}, address)
	spew.Dump(err)

	return t.Int64()
}

func authTx(client *ethclient.Client) (*bind.TransactOpts, error) {
	prikeyByte, err := os.ReadFile("prikey.txt")
	prikey := strings.Trim(string(prikeyByte), "\n")

	privateKey, err := crypto.HexToECDSA(prikey)
	if err != nil {
		panic(err)
	}

	// 公鑰
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println("fromAddress", fromAddress)

	// 產生隨機數
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println("nounce=", nonce)
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		panic(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		panic(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(3000000) // in units

	return auth, nil
}

// 造幣
func (a *amtSrc) Mint(toAddr string, valueData string, nonce, sigHex string) error {

	to := common.HexToAddress(toAddr)
	value := new(big.Int)
	value.SetString(valueData, 10)

	tx, err := authTx(a.client.Client)
	gasPrice, err := a.client.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	tx.GasPrice = gasPrice
	transcation, err := a.token.Mint(tx, to, value)
	if err != nil {
		spew.Dump(err)
		spew.Dump(transcation)
		log.Fatal(err)
		return err
	}

	log.Printf("transcation hash : %s", transcation.Hash())
	return nil

}

// 授權轉移數量
func (a *amtSrc) Approve(senderAddr, spenderAddr string, valueData string) error {

	sender := common.HexToAddress(senderAddr)
	spender := common.HexToAddress(spenderAddr)
	value := new(big.Int)
	value.SetString(valueData, 10)
	transaction, err := a.token.Approve(&bind.TransactOpts{
		From: sender,
	}, spender, value)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(transaction.Hash())
	return nil
}

// 提領授權的代幣
func (a *amtSrc) TransferFrom(senderAddr, fromAddr, toAddr string, valueData, nonce, sigHex string) error {

	from := common.HexToAddress(fromAddr)
	to := common.HexToAddress(toAddr)
	value := new(big.Int)
	value.SetString(valueData, 10)

	tx, err := authTx(a.client.Client)
	transaction, err := a.token.TransferFrom(tx, from, to, value)

	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println(transaction.Hash())
	return nil

}

// 轉移代幣
func (a *amtSrc) Transfer(senderAddr, toAddr string, valueData, nonce, sigHex string) error {

	to := common.HexToAddress(toAddr)
	value := new(big.Int)
	value.SetString(valueData, 10)

	tx, err := authTx(a.client.Client)

	gasPrice, err := a.client.Client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	tx.GasPrice = gasPrice
	transcation, err := a.token.Transfer(tx, to, value)
	if err != nil {
		spew.Dump(err)
		spew.Dump(transcation)
		log.Fatal(err)
		return err
	}

	log.Println(transcation.Hash())
	return nil
}
