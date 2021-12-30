package blockchain

import (
	"reflect"
	"sync"
	"testing"

	"github.com/hskim881028/goblockchain/utility"
)

type fakeDB struct {
	fakeGetBlock func() []byte
	fakeGetChain func() []byte
}

func (f fakeDB) GetBlock(hash string) []byte {
	return f.fakeGetBlock()
}
func (f fakeDB) GetChain() []byte {
	return f.fakeGetChain()
}
func (fakeDB) PutBlock(hash string, data []byte) {}
func (fakeDB) DeleteAllBlocks()                  {}
func (fakeDB) PutChain(data []byte)              {}

func TestReplace(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "xx",
	}
	blocks := []*Block{
		{Difficulty: 2, Hash: "new_test"},
		{Difficulty: 2, Hash: "new_test"},
		{Difficulty: 2, Hash: "new_test"},
	}
	bc.Replace(blocks)
	if bc.CurrentDifficulty != 2 || bc.Height != 3 || bc.NewestHash != "new_test" {
		t.Error("Replace should mutate the blockchain")
	}
}

func TestAddPeerBlock(t *testing.T) {
	bc := &blockchain{
		Height:            1,
		CurrentDifficulty: 1,
		NewestHash:        "test",
	}
	m.Txs["test"] = &Tx{}
	nb := &Block{
		Difficulty: 2,
		Hash:       "new_test",
		Transactions: []*Tx{
			{ID: "test"},
		},
	}
	bc.AddPeerBlock(nb)
	if bc.CurrentDifficulty != 2 || bc.Height != 2 || bc.NewestHash != "new_test" {
		t.Error("AddPeerBlock should mutate the blockchain")
	}
}
func TestGetDifficulty(t *testing.T) {
	blocks := []*Block{
		{PrevHash: "e", Height: 5},
		{PrevHash: "d", Height: 4},
		{PrevHash: "c", Height: 3},
		{PrevHash: "b", Height: 2},
		{PrevHash: "a", Height: 1},
		{PrevHash: "", Height: 0},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeGetBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utility.ToBytes(blocks[fakeBlock])
		},
	}
	type test struct {
		height int
		want   int
	}

	tests := []test{
		{0, defaultDifficulty},
		{2, defaultDifficulty},
		{5, 3},
	}

	for _, tc := range tests {
		bc := &blockchain{Height: tc.height, CurrentDifficulty: defaultDifficulty}
		get := getDifficulty(bc)
		if get != tc.want {
			t.Errorf("getDifficulty should return %d get %d", tc.want, get)
		}
	}
}

func TestFindTx(t *testing.T) {
	t.Run("Tx should not found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeGetBlock: func() []byte {
				b := &Block{
					Height:       1,
					Transactions: []*Tx{},
				}
				return utility.ToBytes(b)
			},
		}
		bc := &blockchain{}
		tx := FindTx(bc, "test")
		if tx != nil {
			t.Error("Tx should not found")
		}
	})

	t.Run("Tx should found", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeGetBlock: func() []byte {
				b := &Block{
					Height: 1,
					Transactions: []*Tx{
						{ID: "test_id"},
					},
				}
				return utility.ToBytes(b)
			},
		}
		bc := &blockchain{}
		tx := FindTx(bc, "test_id")
		if tx == nil {
			t.Error("Tx should found")
		}
	})
}

func TestGetBlocks(t *testing.T) {
	blocks := []*Block{
		{PrevHash: "a"},
		{PrevHash: ""},
	}
	fakeBlock := 0
	dbStorage = fakeDB{
		fakeGetBlock: func() []byte {
			defer func() {
				fakeBlock++
			}()
			return utility.ToBytes(blocks[fakeBlock])
		},
	}
	bc := &blockchain{}
	blockResult := GetBlocks(bc)
	if reflect.TypeOf(blockResult) != reflect.TypeOf([]*Block{}) {
		t.Error("GetBlocks should return a slice of blocks")
	}
}

func TestBlockChain(t *testing.T) {
	t.Run("Should create blockchain", func(t *testing.T) {
		dbStorage = fakeDB{
			fakeGetChain: func() []byte {
				return nil
			},
		}
		bc := Blockchain()
		if bc.Height != 1 {
			t.Error("Blockchain should create a blockchain")
		}
	})

	t.Run("Should restore blockchain", func(t *testing.T) {
		chainOnce = *new(sync.Once)
		dbStorage = fakeDB{
			fakeGetChain: func() []byte {
				bc := blockchain{Height: 2, NewestHash: "test", CurrentDifficulty: 1}
				return utility.ToBytes(bc)
			},
		}
		bc := Blockchain()
		if bc.Height != 2 {
			t.Error("Blockchain should restore a blockchain")
		}
	})
}
