package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
)

type Block struct {
	PrevBlock [sha256.Size]byte
	Author    [32]byte
	Message   []byte
	Nonce     uint32
}

type Tx struct {
	From   [32]byte
	To     [32]byte
	Amount int
}

func getTx(b Block) Tx {
	tx, success := sign.Open(nil, b.Message, &b.Author)
	var result Tx
	if success {
		json.Unmarshal([]byte(tx), &result)
		return result
	}
	return result
}

func hash(block *Block) [sha256.Size]byte {
	blockbytes, _ := json.Marshal(block)
	return sha256.Sum256(blockbytes)
}

func mine(block *Block, target [sha256.Size]byte) bool {
	var candidate uint32
	for {
		block.Nonce = candidate
		hashBlock := hash(block)
		var bHash []byte = hashBlock[:]
		if (bytes.Compare(bHash, target[:])) < 1 {
			fmt.Println(bHash)
			return true
		}
		candidate++
		if candidate == 0 {
			break
		}
	}
	return false
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

	var prev [32]byte
	var block = Block{PrevBlock: prev, Author: *fromPub, Message: signedMessage, Nonce: 0}
	//fmt.Println(block)
	//fmt.Println(getTx(block) == tx)

	//	fmt.Println(hash(&block))
	var target = [sha256.Size]byte{0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	mine(&block, target)
}
