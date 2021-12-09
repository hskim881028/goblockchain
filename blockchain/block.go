package blockchain

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type Block struct {
	Data       string `json:data`
	Hash       string `json:hash`
	PrevHash   string `json:prevHash,omitempty`
	Height     int    `json:height`
	Difficulty int    `json:difficulty`
	Nonce      int    `json:nonce`
}

var ErrorNotFound = errors.New("block not found")

const difficulty int = 2

func (b *Block) Persist() {
	db.SaveBlock(b.Hash, utility.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utility.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", difficulty)
	for {
		blockAsString := fmt.Sprint(b)
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(blockAsString)))
		fmt.Printf("Block as String:%s\nHash:%s\nTarget:%s\nNonce:%d\n\n\n", blockAsString, hash, target, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(data string, prevHash string, height int) *Block {
	block := &Block{
		Data:       data,
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}
	block.mine()
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
