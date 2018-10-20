package main

import (
	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Kita terhubung yey!")
	_ = client // we'll use this in the next section
}