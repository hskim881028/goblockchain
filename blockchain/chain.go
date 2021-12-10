package blockchain

import (
	"fmt"
	"sync"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type blockchain struct {
	NewestHash        string `json:newestHash`
	Height            int    `json:height`
	CurrentDifficulty int    `json:currentDifficulty`
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var b *blockchain
var once sync.Once

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return b.calculateDifficulty()
	} else {
		return b.CurrentDifficulty
	}
}

func (b *blockchain) calculateDifficulty() int {
	blocks := b.Blocks()
	newestBlcok := blocks[0]
	// 해당 인덱스의 블록이 없을지도
	lastCalculatedBlock := blocks[difficultyInterval-1]
	actualTime := (newestBlcok.Timestamp / 60) - (lastCalculatedBlock.Timestamp / 60)
	expectedTime := difficultyInterval * blockInterval

	if actualTime <= (expectedTime - allowedRange) {
		return b.CurrentDifficulty + 1
	} else if actualTime >= (expectedTime + allowedRange) {
		return b.CurrentDifficulty - 1
	}
	return b.CurrentDifficulty
}

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
	b.CurrentDifficulty = block.Difficulty
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
			b = &blockchain{
				Height: 0,
			}
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
