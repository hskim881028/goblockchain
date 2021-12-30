package blockchain

import (
	"reflect"
	"testing"
)

func TestCreateBlock(t *testing.T) {
	dbStorage = fakeDB{}
	Mempool().Txs["test"] = &Tx{}
	b := createBlock("t", 1, 1)
	if reflect.TypeOf(b) != reflect.TypeOf(&Block{}) {
		t.Error("createBlock should return an instance of block")
	}
}

// func TestGetBlock(t *testing.T) {
// 	t.Run("Block not found", func(t *testing.T) {
// 		dbStorage = fakeDB{
// 			fakeGetBlock: func() []byte {
// 				return nil
// 			},
// 		}
// 		_, err := GetBlock("test")
// 		if err == nil {
// 			t.Error("GetBlock should not found")
// 		}
// 	})

// 	t.Run("Block found", func(t *testing.T) {
// 		dbStorage = fakeDB{
// 			fakeGetBlock: func() []byte {
// 				b := &Block{
// 					Height: 1,
// 				}
// 				return utility.ToBytes(b)
// 			},
// 		}
// 		block, _ := GetBlock("test")
// 		if block.Height != 1 {
// 			t.Error("GetBlock shout be found")
// 		}
// 	})
// }
