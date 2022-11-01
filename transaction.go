package tronWallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/api"
	"github.com/golang/protobuf/proto"
	"tronWallet/enums"
	"tronWallet/grpcClient"
)

func createTransactionInput(node enums.Node, fromAddressBase58 string, toAddressBase58 string, amountInSun int64) (*api.TransactionExtention, error) {

	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return nil, err
	}

	return c.Transfer(fromAddressBase58, toAddressBase58, amountInSun)
}

func signTransaction(transaction *api.TransactionExtention, privateKey *ecdsa.PrivateKey) (*api.TransactionExtention, error) {

	rawData, err := proto.Marshal(transaction.Transaction.GetRawData())
	if err != nil {
		return transaction, fmt.Errorf("proto marshal tx raw data error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return transaction, fmt.Errorf("sign error: %v", err)
	}

	transaction.Transaction.Signature = append(transaction.Transaction.Signature, signature)
	return transaction, nil
}

func broadcastTransaction(node enums.Node, transaction *api.TransactionExtention) error {

	c, err := grpcClient.GetGrpcClient(node)
	if err != nil {
		return err
	}

	res, err := c.Broadcast(transaction.Transaction)
	if err != nil {
		return err
	}

	if res.Result != true {
		return errors.New(res.Code.String())
	}

	return nil
}
