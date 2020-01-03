package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
)

type Block struct {
	PrevBlock []byte
	Author    [32]byte
	Message   []byte
	Nonce     int
}

type Tx struct {
	From   [32]byte
	To     [32]byte
	Amount int
}

func main() {
	fromPub, fromPr, _ := sign.GenerateKey(rand.Reader)
	toPub, _, _ := sign.GenerateKey(rand.Reader)
	var tx = Tx{From: *fromPub, To: *toPub, Amount: 10}
	//fmt.Println(tx)

	txJson, _ := json.Marshal(&tx)

	//var m Tx
	//json.Unmarshal([]byte(txJson), &m)
	//fmt.Println(m)

	signedMessage := sign.Sign(nil, txJson, fromPr)
	message, _ := sign.Open(nil, signedMessage, fromPub)
	var m Tx
	json.Unmarshal([]byte(message), &m)
	//fmt.Println(m)

	var prev []byte
	var block = Block{PrevBlock: prev, Author: *fromPub, Message: signedMessage, Nonce: 0}
	fmt.Println(block)

}
