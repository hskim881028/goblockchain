package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hskim881028/goblockchain/blockchain"
	"github.com/hskim881028/goblockchain/utility"
)

var port string

type url string

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:errorMessage`
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", port, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "GET",
			Description: "Get all blocks",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add a block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{Height}"),
			Method:      "GET",
			Description: "See a block",
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
		var addBlockBody addBlockBody
		utility.HandleErr(json.NewDecoder(r.Body).Decode(&addBlockBody))
		blockchain.GetBlcokchain().AddBlock(addBlockBody.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height, err := strconv.Atoi(vars["height"])
	utility.HandleErr(err)
	block, err := blockchain.GetBlcokchain().GetBlock(height)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorResponse{fmt.Sprint(err)})
	} else {
		encoder.Encode(block)
	}
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on Http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
