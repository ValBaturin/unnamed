package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
)

type Tx struct {
	From   [32]byte
	To     [32]byte
	Amount int
}

func main() {
	fromPub, _, _ := sign.GenerateKey(rand.Reader)
	toPub, _, _ := sign.GenerateKey(rand.Reader)
	var tx = Tx{From: *fromPub, To: *toPub, Amount: 10}
	fmt.Println(tx)

	txJson, _ := json.Marshal(&tx)

	var m Tx
	json.Unmarshal([]byte(txJson), &m)
	fmt.Println(m)

	//signedMessage := sign.Sign(nil, data, fromPr)
	//message, _ := sign.Open(nil, signedMessage, fromPub)
	//fmt.Println(tx)
	//fmt.Println(message)

}
