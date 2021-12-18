package p2p

import (
	"encoding/json"
	"fmt"

	"github.com/hskim881028/goblockchain/blockchain"
	"github.com/hskim881028/goblockchain/utility"
)

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
)

type MessageKind int

type Message struct {
	Kind     MessageKind
	Playload []byte
}

func makeMessage(kind MessageKind, payload interface{}) []byte {
	m := Message{
		Kind:     kind,
		Playload: utility.ToJson(payload),
	}
	return utility.ToJson(m)
}

func sendNewestBlock(p *peer) {
	block, err := blockchain.FindBlock(blockchain.Blcokchain().NewestHash)
	utility.HandleError(err)
	message := makeMessage(MessageNewestBlock, block)
	p.inbox <- message
}

func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		var payload blockchain.Block
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		fmt.Println(payload)
	}
}
