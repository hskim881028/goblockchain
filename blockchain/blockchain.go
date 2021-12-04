package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

type Block struct {
	Data     string `json:data`
	Hash     string `json:hash`
	PrevHash string `json:prevHash,omitempty`
	Height   int    `json:height`
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once
var ErrNotFound = errors.New("block not found")

func (b *Block) calculateHash() {
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

func createBlock(data string) *Block {
	newBlock := Block{data, "", getLashHash(), len(GetBlcokchain().blocks) + 1}
	newBlock.calculateHash()
	return &newBlock
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockchain) GetBlock(height int) (*Block, error) {
	if height > len(b.blocks) {
		return nil, ErrNotFound
	}

	return b.blocks[height-1], nil
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

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}
