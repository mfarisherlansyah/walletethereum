package main

import (
	"context"
	// "crypto/ecdsa"
	"fmt"
	"encoding/hex"
	"log"
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	// gocrypto "github.com/alivanz/go-crypto"
	"github.com/alivanz/go-crypto/bitcoin"
	"github.com/btcsuite/btcd/btcec"
	"errors"
)


func EtherSignatureFromECDSA(r,s *big.Int, privkey []byte) ([]byte,error){
	curve := btcec.S256()
	sig, err := SignatureCompact(r,s, curve, false)
	if err != nil {
		return nil, err
	}

	// Convert to Ethereum signature format with 'recovery id' v at the end.
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

/*

func (wallet *ethwallet) Sign(hash []byte) ([]byte,error){
	sig, err := crypto.Sign(hash, wallet.private)
	if err != nil {
		return nil, err
	}
	return sig, nil
} */

/* func TxMake(from string, to string) ([]byte, error) {

	client, err := ethclient.Dial("https://rinkeby.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	var data []byte
	nonce, err := client.PendingNonceAt(context.Background(), from)
	if err != nil {
		log.Fatal(err)
	}
	txs := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, data)
	return txs
} */


/* func Sign(hash []byte) ([]byte, error){
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	return sig, nil
} */

func main() {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyBytes, _ := hex.DecodeString("3f41ea069dcb4b50d73b4928a489f2f8dc3148f6873d00d141856e3f52ecccab")
	wallet, _ := bitcoin.NewWallet(privateKeyBytes)
	pubkey, _ := wallet.PubKey()

	// publicKeyBytes := crypto.FromECDSAPub(&pubkey)

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	/*if !ok {
		log.Fatal("error casting public key to ECDSA")
	} */

	fromAddress := crypto.PubkeyToAddress(pubkey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	gasLimit := uint64(21000) // in units
	
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// toAddress := common.HexToAddress("0x4592d8f8d7b001e72cb26a73e4fa1806a51ac79d")
	toAddress := common.HexToAddress("0x293f51bcc3cf7e73a21c90a20bf9581009047013")
	
	var data []byte
	
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	// tx := gocrypto.TxMake(fromAddress, toAddress)
	// signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	// signedTx, err := wallet.Sign(tx, privateKey)
	/* r,s, err := wallet.Sign(hash.Bytes())
	//Sign(tx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	} */

	chainID, err := client.NetworkID(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    signedTx, err := SignTx(tx, types.NewEIP155Signer(chainID), privateKeyBytes)
    if err != nil {
        log.Fatal(err)
    }

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func SignTx(tx *types.Transaction, signer types.Signer, prv []byte) (*types.Transaction, error) {
	h := signer.Hash(tx) 

	wallet, _ := bitcoin.NewWallet(prv)

	r,s, err := wallet.Sign(h.Bytes())
	if err != nil {
		return nil, err
	}

	signature, err := EtherSignatureFromECDSA(r,s, prv)
	if err != nil {
		log.Fatal(err)
	}

	return tx.WithSignature(signer, signature)
}

func SignatureCompact(r,s *big.Int, curve *btcec.KoblitzCurve, isCompressedKey bool) ([]byte,error){
	sig := btcec.Signature{R:r, S:s}
	// bitcoind checks the bit length of R and S here. The ecdsa signature
	// algorithm returns R and S mod N therefore they will be the bitsize of
	// the curve, and thus correctly sized.

	for i := 0; i < (curve.H+1)*2; i++ {
		// pk, err := recoverKeyFromSignature(curve, sig, hash, i, true)
		// if pk.X.Cmp(key.X) == 0 && pk.Y.Cmp(key.Y) == 0 {
			result := make([]byte, 1, 2*(curve.BitSize/8)+1)
			result[0] = 27 + byte(i)
			if isCompressedKey {
				result[0] += 4
			}
			// Not sure this needs rounding but safer to do so.
			curvelen := (curve.BitSize + 7) / 8

			// Pad R and S to curvelen if needed.
			bytelen := (sig.R.BitLen() + 7) / 8
			if bytelen < curvelen {
				result = append(result,
					make([]byte, curvelen-bytelen)...)
			}
			result = append(result, sig.R.Bytes()...)

			bytelen = (sig.S.BitLen() + 7) / 8
			if bytelen < curvelen {
				result = append(result,
					make([]byte, curvelen-bytelen)...)
			}
			result = append(result, sig.S.Bytes()...)

			return result, nil
		}
	 
	return nil, errors.New("no valid solution for pubkey found")
}