package main

import (
	"context"
	"encoding/hex"
	"log"
	"math/big"
	"os"
	"signer/abis"
	"signer/utils"
	"signer/wallet"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/urfave/cli/v2"
)

func sign(c *cli.Context) error {
	log.Println("signing...")

	wethAddress := common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
	wsbAddress := common.HexToAddress(c.String("WSBAddress"))
	unscaledWethAmount := big.NewInt(17) // 17 ETH ~ 30k
	wethAmount := utils.EthToWei(decimal.NewFromInt(unscaledWethAmount.Int64()))
	unscaledWsbAmount := big.NewInt(41652000000) // 41,652,000,000 WSB  == 60% of 69,420,000,000
	wsbAmount := utils.EthToWei(decimal.NewFromInt(unscaledWsbAmount.Int64()))

	// load vars
	wallet_path := c.String("walletsPath")
	pw_path := c.String("pwPath")
	pub := c.String("publicKey")
	wallet_addr := common.HexToAddress(pub)

	// import private key
	wallet.ImportPKey(c)

	// Open wallet
	ew, err := wallet.NewExecWallet(wallet_path, pw_path)
	if err != nil {
		log.Println("Failed opening wallet")
		return err
	}

	// Load key
	rpc, err := rpc.Dial(c.String("rpcAddress"))
	if err != nil {
		log.Println("Error connecting to network: ", err)
	}
	client := ethclient.NewClient(rpc)

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	log.Println("Connected to ethereum client, latest block:", header.Number.String())

	cachedKey := ew.GetAccountByAddr(client, wallet_addr)
	bal, err := client.BalanceAt(context.Background(), cachedKey.Account.Address, nil)
	if err != nil {
		log.Println("Error getting account ETH balance: ", err)
	}

	log.Println("-- Account Loaded --")
	log.Println("Address:", cachedKey.Account.Address.Hex())
	log.Println("Nonce:", cachedKey.Nonce.String())
	log.Println("Balance (ETH):", utils.WeiToEth(bal).String())
	log.Println("Continue if this looks about right")
	if !utils.GetUserConfirmation() {
		log.Println("Exiting on user request...")
	}

	v2Transactor, err := abis.NewIUniswapV2Router02(common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"), client)
	if err != nil {
		log.Println("Error loading contract: ", err)
	}

	auth1, err := bind.NewKeyedTransactorWithChainID(cachedKey.Key.PrivateKey, big.NewInt(1))
	if err != nil {
		log.Printf("rawContractCall: %s\n", err)
		return err
	}

	auth1.Nonce = big.NewInt(int64(cachedKey.Nonce.Uint64()))

	//
	// Gas Calculation
	suggestedFeeCap, _ := client.SuggestGasPrice(context.Background())
	suggestedTipCap, _ := client.SuggestGasTipCap(context.Background())
	auth1.GasFeeCap = suggestedFeeCap // 1559
	auth1.GasTipCap = suggestedTipCap // 1559
	auth1.GasLimit = 1_000_000
	auth1.NoSend = true

	// sign tx
	tx, err := v2Transactor.AddLiquidity(auth1, wethAddress, wsbAddress, wethAmount, wsbAmount, wethAmount, wsbAmount, cachedKey.Account.Address, big.NewInt(999999999999))
	if err != nil {
		log.Println("Error signing tx: ", err)
	}

	binary, err := tx.MarshalBinary()
	if err != nil {
		log.Println("Error marshalling tx: ", err)
	}
	signedTxString := hex.EncodeToString(binary)
	log.Println("Successfully signed the transaction!")

	log.Println("###########################################")
	log.Println(signedTxString)
	log.Println("###########################################")
	return nil
}
func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	app := &cli.App{
		Name:  "signer",
		Usage: "",
	}
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "sign",
			Action: sign,
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "rpcAddress", Value: "https://rpc.ankr.com/eth"},
				&cli.StringFlag{Name: "WSBAddress", Value: "TODO"},
				&cli.StringFlag{Name: "publicKey", Value: "TODO"},
				&cli.StringFlag{Name: "privateKey", Value: "TODO"},
				&cli.StringFlag{Name: "walletsPath", Value: "./wallets/"},
				&cli.StringFlag{Name: "pwPath", Value: "pw.ptxt"},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
