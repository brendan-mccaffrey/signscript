package wallet

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Address struct {
	Addr string `json:"address"`
}

type Addresses struct {
	Addrs []*Address `json:"bots"`
}

func GetBalances(client *ethclient.Client, fromAddr common.Address, totalAmt *big.Int, accounts []common.Address, tokenAddr string) []*big.Int {
	var balances []*big.Int

	// append each wallet balance
	for i := 0; i < len(accounts); i++ {
		bal, err := client.BalanceAt(context.Background(), accounts[i], nil)
		if err != nil {
			log.Println("balanceAt ", err)
			return make([]*big.Int, 0)
		}
		balances = append(balances, bal)
	}

	return balances
}
