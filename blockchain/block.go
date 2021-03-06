package blockchain

import (
	"errors"
	"strings"
	"time"

	"github.com/hskim881028/goblockchain/utility"
)

type Block struct {
	Hash         string `json:"hash"`
	PrevHash     string `json:"prevHash,omitempty"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
	Transactions []*Tx  `json:"transactions"`
}

var ErrorNotFound = errors.New("block not found")

func persistBlock(b *Block) {
	dbStorage.PutBlock(b.Hash, utility.ToBytes(b))
}

func (b *Block) restore(data []byte) {
	utility.FromBytes(b, data)
}

func (b *Block) mine() {
	target := strings.Repeat("0", b.Difficulty)
	for {
		b.Timestamp = int(time.Now().Unix())
		hash := utility.Hash(b)
		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			break
		} else {
			b.Nonce++
		}
	}
}

func createBlock(prevHash string, height, difficulty int) *Block {
	block := &Block{
		Hash:       "",
		PrevHash:   prevHash,
		Height:     height,
		Difficulty: difficulty,
		Nonce:      0,
	}
	block.Transactions = Mempool().TxToConfirm()
	block.mine()
	persistBlock(block)
	return block
}

func GetBlock(hash string) (*Block, error) {
	blockBytes := dbStorage.GetBlock(hash)
	if blockBytes == nil {
		return nil, ErrorNotFound
	}

	block := &Block{}
	block.restore(blockBytes)
	return block, nil
}
