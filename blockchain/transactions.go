package blockchain

import (
	"errors"
	"fmt"
	"time"

	"github.com/hskim881028/goblockchain/utility"
)

const (
	mineReward int = 50
)

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"txIns"`
	TxOuts    []*TxOut `json:"txOuts"`
}

type TxIn struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type mempool struct {
	Txs []*Tx
}

var Mempool *mempool = &mempool{}

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

func makeTx(from, to string, amount int) (*Tx, error) {
	if Blcokchain().BalanceByAddress(from) < amount {
		err := fmt.Sprintf("[not enough money] from : %s", from)
		return nil, errors.New(err)
	}

	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	preTxOuts := Blcokchain().TxOutsByAddress(from)
	for _, txOut := range preTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{txOut.Owner, txOut.Amount}
		txIns = append(txIns, txIn)
		total += txOut.Amount
	}

	change := total - amount
	if change != 0 {
		changeTxOut := &TxOut{from, change}
		txOuts = append(txOuts, changeTxOut)
	}
	txOut := &TxOut{
		Owner:  to,
		Amount: amount,
	}
	txOuts = append(txOuts, txOut)

	tx := &Tx{"", int(time.Now().Unix()), txIns, txOuts}
	tx.getId()
	return tx, nil
}

func (m *mempool) AddTx(to string, amount int) error {
	tx, err := makeTx("hskim", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}
