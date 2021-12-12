package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hskim881028/goblockchain/blockchain"
	"github.com/hskim881028/goblockchain/utility"
	"github.com/hskim881028/goblockchain/wallet"
)

var port string

type url string

type addBlockBody struct {
	Message string
}

type errorResponse struct {
	ErrorMessage string `json:errorMessage`
}

type balanceResponse struct {
	Address string `json:"address"`
	Balance int    `json:"balance`
}

type myWalletResponse struct {
	Address string `json:"address"`
}

type addTxPayload struct {
	To     string
	Amount int
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
		{
			URL:         url("/status"),
			Method:      "GET",
			Description: "Status of blockchain",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an Address",
		},
	}
	utility.HandleError(json.NewEncoder(rw).Encode(data))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		utility.HandleError(json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.Blcokchain())))
		return
	case "POST":
		blockchain.Blcokchain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
		return
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrorNotFound {
		utility.HandleError(encoder.Encode(errorResponse{fmt.Sprint(err)}))
	} else {
		utility.HandleError(encoder.Encode(block))
	}
}

func status(rw http.ResponseWriter, r *http.Request) {
	utility.HandleError(json.NewEncoder(rw).Encode(blockchain.Blcokchain()))
}

func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	total := r.URL.Query().Get("total")

	switch total {
	case "true":
		amount := blockchain.BalanceByAddress(blockchain.Blcokchain(), address)
		json.NewEncoder(rw).Encode(balanceResponse{address, amount})
	default:
		utility.HandleError(json.NewEncoder(rw).Encode(blockchain.UTxOutsByAddress(blockchain.Blcokchain(), address)))
	}
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	utility.HandleError(json.NewEncoder(rw).Encode(blockchain.Mempool.Txs))
}

func transactions(rw http.ResponseWriter, r *http.Request) {
	var payload addTxPayload
	utility.HandleError(json.NewDecoder(r.Body).Decode(&payload))
	err := blockchain.Mempool.AddTx(payload.To, payload.Amount)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(errorResponse{err.Error()})
		return
	}
	rw.WriteHeader(http.StatusCreated)
}

func myWallet(rw http.ResponseWriter, r *http.Request) {
	address := wallet.Wallet().Address
	json.NewEncoder(rw).Encode(myWalletResponse{Address: address})
}

func Start(aPort int) {
	port = fmt.Sprintf(":%d", aPort)
	router := mux.NewRouter()
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transactions", transactions).Methods("POST")
	router.HandleFunc("/wallet", myWallet).Methods("GET")
	fmt.Printf("Listening on Http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, router))
}
