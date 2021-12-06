package blockchain

import (
	"crypto/sha256"
	"fmt"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type Block struct {
	Data     string `json:data`
	Hash     string `json:hash`
	PrevHash string `json:prevHash,omitempty`
	Height   int    `json:height`
}

func (b *Block) Persist() {
	db.SaveBlock(b.Hash, utility.ToBytes(b))
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:     data,
		Hash:     "",
		PrevHash: prevHash,
		Height:   height,
	}
	payload := block.Data + block.PrevHash + fmt.Sprint(block.Height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	block.Persist()
	return block
}
