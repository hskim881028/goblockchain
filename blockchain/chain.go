package blockchain

import (
	"fmt"
	"sync"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type blockchain struct {
	NewestHash string `json:"newestHash`
	Height     int    `json:height`
}

var b *blockchain
var once sync.Once

func (b *blockchain) persist() {
	db.SaveBlockChain(utility.ToBytes(b))
}

func (b *blockchain) restore(data []byte) {
	utility.FromBytes(b, data)
}

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash

	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Blocks() []*Block {
	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := FindBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}
func Blcokchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			checkPoint := db.CheckPoint()
			if checkPoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.restore(checkPoint)
			}
		})
	}
	fmt.Println(b.NewestHash)
	return b
}
