package ethereum

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthService struct {
	Client *ethclient.Client
}

func New() *EthService {
	log.Println("Connent to Ethereum...")
	client, err := ethclient.Dial("https://ethereum-goerli-rpc.allthatnode.com")
	if err != nil {
		log.Fatal("New EthService fail")
		return nil
	}
	return &EthService{
		Client: client,
	}
}
