package utils

import (
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
)

func GetUserConfirmation() bool {
	log.Print("Would you like to continue [y/n]? ")
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "y" {
		return false
	}
	return true
}

// WeiToEth automatically inputs 18 decimals
func WeiToEth(iamount *big.Int) decimal.Decimal {
	return ToDecimal(iamount, 18)
}

// EthToWei automatically inputs 18 decimals
func EthToWei(iamount decimal.Decimal) *big.Int {
	return ToWei(iamount, 18)
}

// ToWei decimals to wei
func ToWei(iamount interface{}, decimals int) *big.Int {
	amount := decimal.NewFromFloat(0)
	switch v := iamount.(type) {
	case string:
		amount, _ = decimal.NewFromString(v)
	case float64:
		amount = decimal.NewFromFloat(v)
	case int64:
		amount = decimal.NewFromFloat(float64(v))
	case decimal.Decimal:
		amount = v
	case *decimal.Decimal:
		amount = *v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	result := amount.Mul(mul)

	wei := new(big.Int)
	wei.SetString(result.String(), 10)

	return wei
}

func TxFrom(tx *types.Transaction) common.Address {
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		from, err = types.Sender(types.HomesteadSigner{}, tx)
		if err != nil {
			log.Println("Utils: Error getting sender from tx", err)
			return common.HexToAddress("0x0000000000000000000000000000000000000000")
		}
	}
	return from
}

// ToDecimal wei to decimals
func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}
