package blockchain

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/hskim881028/goblockchain/db"
	"github.com/hskim881028/goblockchain/utility"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"currentDifficulty"`
	m                 sync.Mutex
}

type storage interface {
	GetBlock(hash string) []byte
	PutBlock(hash string, data []byte)
	DeleteAllBlocks()
	GetChain() []byte
	PutChain(data []byte)
}

const (
	defaultDifficulty  int = 2
	difficultyInterval int = 5
	blockInterval      int = 2
	allowedRange       int = 2
)

var b *blockchain
var chainOnce sync.Once
var dbStorage storage = db.DB{}

func (b *blockchain) restore(data []byte) {
	utility.FromBytes(b, data)
}

func (b *blockchain) AddBlock() *Block {
	block := createBlock(b.NewestHash, b.Height+1, getDifficulty(b))
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	persistBlockChain(b)
	return block
}

func (b *blockchain) Replace(newBlocks []*Block) {
	b.m.Lock()
	defer b.m.Unlock()

	b.CurrentDifficulty = newBlocks[0].Difficulty
	b.NewestHash = newBlocks[0].Hash
	b.Height = len(newBlocks)

	persistBlockChain(b)
	dbStorage.DeleteAllBlocks()
	for _, block := range newBlocks {
		persistBlock(block)
	}
}

func (b *blockchain) AddPeerBlock(newBlock *Block) {
	b.m.Lock()
	m.m.Lock()
	defer b.m.Unlock()
	defer m.m.Unlock()

	b.Height += 1
	b.NewestHash = newBlock.Hash
	b.CurrentDifficulty = newBlock.Difficulty

	persistBlockChain(b)
	persistBlock(newBlock)
	for _, tx := range newBlock.Transactions {
		_, ok := m.Txs[tx.ID]
		if ok {
			delete(m.Txs, tx.ID)
		}
	}
}

func persistBlockChain(b *blockchain) {
	dbStorage.PutChain(utility.ToBytes(b))
}

func getDifficulty(b *blockchain) int {
	if b.Height == 0 {
		return defaultDifficulty
	} else if b.Height%difficultyInterval == 0 {
		return calculateDifficulty(b)
	} else {
		return b.CurrentDifficulty
	}
}

func calculateDifficulty(b *blockchain) int {
	blocks := GetBlocks(b)
	newestBlcok := blocks[0]
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

func Txs(b *blockchain) []*Tx {
	var txs []*Tx
	for _, block := range GetBlocks(b) {
		txs = append(txs, block.Transactions...)
	}
	return txs
}

func FindTx(b *blockchain, targetID string) *Tx {
	for _, tx := range Txs(b) {
		if tx.ID == targetID {
			return tx
		}
	}
	return nil
}

func UTxOutsByAddress(b *blockchain, address string) []*UTxOut {
	var uTxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	for _, block := range GetBlocks(b) {
		for _, tx := range block.Transactions {
			for _, input := range tx.TxIns {
				if input.Signature == "COINBASE" {
					break
				}

				if FindTx(b, input.TxID).TxOuts[input.Index].Address == address {
					creatorTxs[input.TxID] = true
				}
			}

			for index, output := range tx.TxOuts {
				if output.Address == address {
					if _, ok := creatorTxs[tx.ID]; !ok {
						uTxOut := &UTxOut{tx.ID, index, output.Amount}
						if !isOnMempool(uTxOut) {
							uTxOuts = append(uTxOuts, uTxOut)
						}
					}
				}
			}
		}
	}

	return uTxOuts
}

func BalanceByAddress(b *blockchain, address string) int {
	txOuts := UTxOutsByAddress(b, address)
	var amount int
	for _, txOut := range txOuts {
		amount += txOut.Amount
	}
	return amount
}

func Status(b *blockchain, rw http.ResponseWriter) {
	b.m.Lock()
	defer b.m.Unlock()
	utility.HandleError(json.NewEncoder(rw).Encode(b))
}

func GetBlocks(b *blockchain) []*Block {
	b.m.Lock()
	defer b.m.Unlock()

	var blocks []*Block
	hashCursor := b.NewestHash
	for {
		block, _ := GetBlock(hashCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			hashCursor = block.PrevHash
		} else {
			break
		}
	}
	return blocks
}

func Blockchain() *blockchain {
	chainOnce.Do(func() {
		b = &blockchain{
			Height: 0,
		}
		checkPoint := dbStorage.GetChain()
		if checkPoint == nil {
			b.AddBlock()
		} else {
			b.restore(checkPoint)
		}
	})
	return b
}
