package main

import (
	"fmt"

	"github.com/hskim881028/goblockchain/blockchain"
)

func main() {
	chain := blockchain.GetBlcokchain()
	chain.AddBlock("Second Block")
	chain.AddBlock("Third Block")
	chain.AddBlock("Fourth Block")
	for _, block := range chain.AllBlocks() {
		fmt.Println(block.Data)
		fmt.Println(block.Hash)
		fmt.Println(block.PrevHash)
	}
}
