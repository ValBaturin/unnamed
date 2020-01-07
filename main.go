package main

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/nacl/sign"
	"math/rand"
)

type Hash [sha256.Size]byte

var zeroHash = Hash{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type Storage struct {
	top  Hash
	data map[Hash]Block
}

func initLedger(s *Storage) map[[32]byte]int {
	current := s.top
	balance := make(map[[32]byte]int)
	for current != zeroHash {
		tx := (s.data[current]).getTx()
		balance[tx.From] -= tx.Amount
		balance[tx.To] += tx.Amount

		current = s.data[current].PrevBlock
	}
	return balance
}

func (s *Storage) init() {
	s.top = zeroHash
	s.data = make(map[Hash]Block)
}

func (s *Storage) add(b *Block) {
	s.data[b.hash()] = *b
	s.top = b.hash()
}

type Block struct {
	PrevBlock Hash
	Tx        SignedTx
	Nonce     uint32
}

type SignedTx struct {
	Message []byte
	Author  [32]byte
}

type Tx struct {
	From   [32]byte
	To     [32]byte
	Amount int
}

func (b Block) getTx() Tx {
	tx, success := sign.Open(nil, b.Tx.Message, &b.Tx.Author)
	var result Tx
	if success {
		json.Unmarshal([]byte(tx), &result)
		return result
	}
	return result
}

func (block *Block) hash() Hash {
	blockbytes, _ := json.Marshal(block)
	return sha256.Sum256(blockbytes)
}

func mine(block *Block, target Hash) bool {
	var candidate uint32
	for {
		block.Nonce = candidate
		hashBlock := block.hash()
		var bHash []byte = hashBlock[:]
		if (bytes.Compare(bHash, target[:])) < 1 {
			//fmt.Println(bHash)
			return true
		}
		candidate++
		if candidate == 0 {
			break
		}
	}
	return false
}

type Account struct {
	Pubkey [32]byte
	Prkey  [64]byte
}

func generateAccounts(n int) []Account {
	var accounts []Account
	for i := 0; i < n; i++ {
		pub, pr, _ := sign.GenerateKey(crand.Reader)
		accounts = append(accounts, Account{Pubkey: *pub, Prkey: *pr})
	}
	return accounts
}

func generateChain(target Hash, length int, acc []Account) *Storage {
	var storage Storage
	storage.init()
	prev := zeroHash

	for i := 0; i < length; i++ {
		sender := acc[rand.Intn(len(acc))]
		receiver := acc[rand.Intn(len(acc))]

		// make receiver different from sender
		for {
			if receiver != sender {
				break
			}
			receiver = acc[rand.Intn(len(acc))]
		}
		amount := rand.Intn(9) + 1 // a single digit coin from 1 to 9

		fmt.Println("chosen sender", sender.Pubkey)
		fmt.Println("chosen receiver", receiver.Pubkey)
		fmt.Println("chosen amount", amount)

		tx := Tx{From: sender.Pubkey, To: receiver.Pubkey, Amount: amount}
		txJson, _ := json.Marshal(&tx)
		signedMessage := sign.Sign(nil, txJson, &sender.Prkey)

		block := Block{PrevBlock: prev, Tx: SignedTx{Author: sender.Pubkey, Message: signedMessage}, Nonce: 0}
		mine(&block, target)
		storage.add(&block)
	}

	return &storage
}

func main() {
	var target = Hash{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	accounts := generateAccounts(2)
	storage := generateChain(target, 2, accounts)
	balance := initLedger(storage)
	fmt.Println(accounts)
	fmt.Println(balance)
}
