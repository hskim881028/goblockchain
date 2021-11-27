package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", hash)
}

func getLashHash() string {
	totalBlocks := len(GetBlcokchain().blocks)
	if totalBlocks == 0 {
		return ""
	}

	return GetBlcokchain().blocks[totalBlocks-1].Hash
}

func createBlock(data string) *block {
	newBlock := block{Data: data, Hash: "", PrevHash: getLashHash()}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func GetBlcokchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("Genesis Block")
		})
	}

	return b
}

func (b *blockchain) AllBlocks() []*block {
	return b.blocks
}
