# Tx Signer

Before this works, you will have to allow the UniswapRouter to spend your tokens. You can do this on the UI by clicking "Approve" on the token you want to spend.

[here](https://app.uniswap.org/#/add/v2/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/undefined)

Hit approve and sign the message, but do not hit "send" or "confirm".

Step 1) Download Go [tutorial](https://go.dev/doc/install)

Step 2) Clone the repo

`git clone xxxxx`

Step 3) Get dependencies

```bash
cd signscript
go get .

# if there are still errors, cd into each directory and run get again
cd wallet
go get .
cd ..
# ...
```

Step 4) Configure main.go

Enter at line 129

```go
&cli.StringFlag{Name: "WSBAddress", Value: "TODO"},
&cli.StringFlag{Name: "publicKey", Value: "TODO"},
&cli.StringFlag{Name: "privateKey", Value: "TODO"},
```

for example

```go
&cli.StringFlag{Name: "WSBAddress", Value: "0xAobcD1234567890..."},
&cli.StringFlag{Name: "publicKey", Value: "0x1234567890..."},
&cli.StringFlag{Name: "privateKey", Value: "982329...."},
```

Step 5) Run the script

```bash
go run main.go sign
```

Step 6) Copy the output and share over signal

DO NOT PUSH OR SHARE ANY CODE AFTER ENTERING YOUR PRIVATE KEY

