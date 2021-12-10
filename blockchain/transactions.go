package blockchain

import (
	"time"

	"github.com/hskim881028/goblockchain/utility"
)

const (
	mineReward int = 50
)

type Tx struct {
	Id        string   `json:id`
	Timestamp int      `json:timestamp`
	TxIns     []*TxIn  `json:txIns`
	TxOuts    []*TxOut `json:txOuts`
}

type TxIn struct {
	Owner  string `json:owner`
	Amount int    `json:amount`
}

type TxOut struct {
	Owner  string `json:owner`
	Amount int    `json:amount`
}

func (t *Tx) getId() {
	t.Id = utility.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"COINBASE", mineReward},
	}
	txOuts := []*TxOut{
		{address, mineReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOuts,
	}
	tx.getId()
	return &tx
}
