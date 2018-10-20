package main

import "flag"

import (
	// "crypto/ecdsa"
	"fmt"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/alivanz/go-crypto/bitcoin"
)

func main() {
	// privateKey, err := crypto.GenerateKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	privateKeyBytes, _ := hex.DecodeString("3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c")
	// fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19
	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("Error casting public key to ECDSA")
	// }

	// var ip = flag.String("PrivateKey", "3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c", "axaxax")

	wallet, _ := bitcoin.NewWallet(privateKeyBytes)
	pubkey, _ := wallet.PubKey()

	publicKeyBytes := crypto.FromECDSAPub(&pubkey)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:]) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	address := crypto.PubkeyToAddress(pubkey).Hex()
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) // 0x96216849c49358b10257cb55b28ea603c874b05e
}