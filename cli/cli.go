package cli

import (
	"flag"
	"fmt"
	"os"

	explorer "github.com/hskim881028/goblockchain/explorer/templates"
	"github.com/hskim881028/goblockchain/rest"
)

func usage() {
	fmt.Printf("Welcome to go blockchain\n\n")
	fmt.Printf("Plase use the following flags:\n\n")
	fmt.Printf("-port : Set the PORT of the server\n")
	fmt.Printf("-moode : Choose between 'rest' and 'html'\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	port := flag.Int("port", 4000, "Set port of the server")
	mode := flag.String("mode", "rest", "Choose between 'html' and 'rest'")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}

	fmt.Println(*port, *mode)
}
