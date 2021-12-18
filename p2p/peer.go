package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type peer struct {
	conn  *websocket.Conn
	inbox chan []byte
}

var Peers map[string]*peer = make(map[string]*peer)

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{conn, make(chan []byte)}
	go p.read()
	go p.write()
	Peers[key] = p
	return p
}

func (p *peer) read() {
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	for {
		m := <-p.inbox
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}
