package p2p

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type peer struct {
	conn    *websocket.Conn
	inbox   chan []byte
	key     string
	address string
	port    string
}

var Peers map[string]*peer = make(map[string]*peer)

func initPeer(conn *websocket.Conn, address, port string) *peer {
	key := fmt.Sprintf("%s:%s", address, port)
	p := &peer{
		conn:    conn,
		inbox:   make(chan []byte),
		key:     key,
		address: address,
		port:    port,
	}
	go p.read()
	go p.write()
	Peers[key] = p
	return p
}

func (p *peer) read() {
	defer p.close()
	for {
		_, m, err := p.conn.ReadMessage()
		if err != nil {
			break
		}
		fmt.Printf("%s", m)
	}
}

func (p *peer) write() {
	defer p.close()
	for {
		m, ok := <-p.inbox
		if !ok {
			break
		}
		p.conn.WriteMessage(websocket.TextMessage, m)
	}
}

func (p *peer) close() {
	p.conn.Close()
	delete(Peers, p.key)
}
