package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	data     string
	hash     string
	prevHash string
}

type blockchain struct {
	blocks []*block
}

var b *blockchain
var once sync.Once

func (b *block) calculateHash() {
	hash := sha256.Sum256([]byte(b.data + b.prevHash))
	b.hash = fmt.Sprintf("%x", hash)
}

func getLashHash() string {
	totalBlocks := len(GetBlcokchain().blocks)
	if totalBlocks == 0 {
		return ""
	}

	return GetBlcokchain().blocks[totalBlocks-1].hash
}

func createBlock(data string) *block {
	newBlock := block{data: data, hash: "", prevHash: getLashHash()}
	newBlock.calculateHash()
	return &newBlock
}

func GetBlcokchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.blocks = append(b.blocks, createBlock("Genesis Block"))
		})
	}

	return b
}
