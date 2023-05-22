package main

import (
	"context"
	"flag"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"os"
)

func main() {

	nodeURL := os.Getenv("INFURA_URL")
	txnHash := flag.String("txnHash", "", "Transaction hash")
	flag.Parse()
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	txHash := common.HexToHash(*txnHash)
	recoverPublicKey(client, txHash)
}

func getSignature(tx *types.Transaction) []byte {
	v, r, s := tx.RawSignatureValues()
	var recoveryID byte
	if v.Sign() == 0 || v.Sign() == 1 {
		recoveryID = byte(v.Uint64())
	} else {
		recoveryID = byte(v.Uint64() - 27)
	}
	signature := append(r.Bytes(), s.Bytes()...)
	signature = append(signature, recoveryID)
	log.Printf("Signature: 0x%x", signature)
	return signature
}

func getTxHash(tx *types.Transaction) []byte {
	signer := types.NewCancunSigner(tx.ChainId())
	hash := signer.Hash(tx)
	log.Printf("Serialized transaction hash: 0x%x", hash)
	return hash.Bytes()
}

func recoverPublicKey(client *ethclient.Client, txHash common.Hash) {
	tx, _, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatalf("Failed to fetch transaction: %v", err)
	}

	signature := getSignature(tx)
	hash := getTxHash(tx)

	publicKey, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		log.Fatalf("Failed to recover public key: %v", err)
	}
	pubKey, err := crypto.UnmarshalPubkey(publicKey)
	if err != nil {
		log.Fatalf("Failed to unmarshal public key: %v", err)
	}
	log.Printf("Recovered public key: 0x%x", pubKey)
	log.Printf("Recovered address: 0x%x", crypto.PubkeyToAddress(*pubKey))
}
