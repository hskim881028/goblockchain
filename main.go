package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hskim881028/goblockchain/blockchain"
	"github.com/hskim881028/goblockchain/utility"
)

const port string = ":4000"

type URL string

type AddBlockBody struct {
	Message string
}

func (u URL) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type URLDescription struct {
	URL         URL    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []URLDescription{
		{
			URL:         URL("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         URL("/blocks"),
			Method:      "GET",
			Description: "Get all blocks",
		},
		{
			URL:         URL("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
	}

	rw.Header().Add("Content-Type", "application/json")
	// b, err := json.Marshal(data)
	// utility.HandleErr(err)
	// fmt.Printf("%s", b)
	json.NewEncoder(rw).Encode(data)

}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlcokchain().AllBlocks())
	case "POST":
		var addBlockBody AddBlockBody
		utility.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlcokchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}

}

func main() {
	http.HandleFunc("/", documentation)
	http.HandleFunc("/blocks", blocks)
	fmt.Printf("Listening on Http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
