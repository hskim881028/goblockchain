package blockchain

import (
	"crypto/sha256"
	"errors"
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

var ErrorNotFound = errors.New("block not found")

func (b *Block) Persist() {
	db.SaveBlock(b.Hash, utility.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utility.FromBytes(b, data)
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

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrorNotFound
	}

	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
