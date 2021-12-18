package p2p

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hskim881028/goblockchain/utility"
)

var conns []*websocket.Conn
var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	openPort := r.URL.Query().Get("openPort")
	ip := utility.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	conn, err := upgrader.Upgrade(rw, r, nil)
	utility.HandleError(err)
	peer := initPeer(conn, ip, openPort)
	time.Sleep(2 * time.Second)
	peer.inbox <- []byte("Hello from Port 3000!")
}

func AddPeer(address, port, openPort string) {
	url := fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort[1:])
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	utility.HandleError(err)
	peer := initPeer(conn, address, port)
	time.Sleep(3 * time.Second)
	peer.inbox <- []byte("Heelo form Port 4000!")
}
