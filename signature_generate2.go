package main

import (
	"fmt"
	"log"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"github.com/alivanz/go-crypto/bitcoin"
	"github.com/btcsuite/btcd/btcec"
	"math/big"
	// "crypto/ecdsa"
	"errors"
	// "github.com/ethereum/go-ethereum/common/math"
	// "github.com/ethereum/go-ethereum/crypto/secp256k1"
	"bytes"
)

func EtherSignatureFromECDSA(r,s *big.Int, privkey []byte) ([]byte,error){
	curve := btcec.S256()
	// btcecpubkey := btcec.PublicKey{
	// 	Curve: curve,
	// 	X: pubkey.X,
	// 	Y: pubkey.Y,
	// }
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

func main() {
	// privateKey, err := crypto.HexToECDSA("3cd0560f5b27591916c643a0b7aa69d03839380a738d2e912990dcc573715d2c")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	privateKeyBytes, _ := hex.DecodeString("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	wallet, _ := bitcoin.NewWallet(privateKeyBytes)
	pubkey, _ := wallet.PubKey()

	publicKeyBytes := crypto.FromECDSAPub(&pubkey)

	data := []byte("hellohellohello") // nanti diganti pake data hasil raw transaction yang mau di-sign
	hash := crypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	r,s, err := wallet.Sign(hash.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	signature, err := EtherSignatureFromECDSA(r,s, privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hexutil.Encode(signature)) // 0x82fc103b2ee3dcce6d0878e388fa34d74c15151d41d8482a59ed8abab13406121e0735d1be8f7e61cf51f0e2975a03e6264b35100a377f969cab998850177c9201
	fmt.Println(hexutil.Encode(r.Bytes()))
	fmt.Println(hexutil.Encode(s.Bytes()))

	sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	fmt.Println(matches) // true
	fmt.Println(hex.EncodeToString(sigPublicKey))
	fmt.Println(hex.EncodeToString(publicKeyBytes))
	sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	fmt.Println(matches) // true
	fmt.Println(hex.EncodeToString(sigPublicKeyBytes))
	fmt.Println(hex.EncodeToString(publicKeyBytes))

	// signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true
}