package blockchain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions`
}

var ErrorNotFound = errors.New("block not found")

func (b *Block) Persist() {
	db.SaveBlock(b.Hash, utility.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utility.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utility.Hash(b)
		fmt.Printf("\n\n\nTarget:%s\nHash:%s\nNonce:%d\n\n\n", target, hash, b.Nonce)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty(Blcokchain()),
		Nonce:      0,
	}
	block.mine()
	block.Transactions = Mempool.TxToConfirm()
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
