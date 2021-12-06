package blockchain

import (
	"bytes"
	"encoding/gob"
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

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.persist()
}

func (b *blockchain) Restore(data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	utility.HandleErr(decoder.Decode(b))
}

func Blcokchain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
			fmt.Printf("Newest Hash : %s\n", b.NewestHash)
			fmt.Printf("Height : %d\n", b.Height)

			checkPoint := db.CheckPoint()
			if checkPoint == nil {
				b.AddBlock("Genesis Block")
			} else {
				b.Restore(checkPoint)
			}
		})
	}
	fmt.Printf("Newest Hash : %s\n", b.NewestHash)
	fmt.Printf("Height : %d\n", b.Height)
	return b
}
