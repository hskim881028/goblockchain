package main

import (
	"github.com/hskim881028/goblockchain/cli"
	"github.com/hskim881028/goblockchain/db"
)

func main() {
	defer db.Close()
	db.InitDB()
	cli.Start()
}
