package p2p

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hskim881028/goblockchain/blockchain"
	"github.com/hskim881028/goblockchain/utility"
)

const (
	MessageNewestBlock MessageKind = iota
	MessageAllBlocksRequest
	MessageAllBlocksResponse
	MessageNewBlockNotify
	MessageNewTxNotify
	MessageNewPeerNotify
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
	block, err := blockchain.GetBlock(blockchain.Blockchain().NewestHash)
	utility.HandleError(err)
	message := makeMessage(MessageNewestBlock, block)
	p.inbox <- message
}

func requestAllBlocks(p *peer) {
	m := makeMessage(MessageAllBlocksRequest, nil)
	p.inbox <- m
}

func sendAllBlocks(p *peer) {
	b := blockchain.Blocks(blockchain.Blockchain())
	m := makeMessage(MessageAllBlocksResponse, b)
	p.inbox <- m
}

func notifyNewBlock(b *blockchain.Block, p *peer) {
	m := makeMessage(MessageNewBlockNotify, b)
	p.inbox <- m
}

func notifyNewTx(tx *blockchain.Tx, p *peer) {
	m := makeMessage(MessageNewTxNotify, tx)
	p.inbox <- m
}

func notifyNewwPeer(address string, p *peer) {
	m := makeMessage(MessageNewPeerNotify, address)
	p.inbox <- m
}

func handleMessage(m *Message, p *peer) {
	switch m.Kind {
	case MessageNewestBlock:
		fmt.Printf("Sending newest block to %s\n", p.key)
		var payload blockchain.Block
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		b, err := blockchain.GetBlock(blockchain.Blockchain().NewestHash)
		utility.HandleError(err)
		if payload.Height >= b.Height {
			fmt.Printf("Requesting all the blocks from %s", p.key)
			requestAllBlocks(p)
		} else {
			fmt.Printf("Sending newest block to %s\n", p.key)
			sendNewestBlock(p)
		}
	case MessageAllBlocksRequest:
		fmt.Printf("%s want all the blocks\n", p.key)
		sendAllBlocks(p)
	case MessageAllBlocksResponse:
		fmt.Printf("Received all the blocks form %s", p.key)
		var payload []*blockchain.Block
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		blockchain.Blockchain().Replace(payload)
	case MessageNewBlockNotify:
		var payload *blockchain.Block
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		blockchain.Blockchain().AddPeerBlock(payload)
	case MessageNewTxNotify:
		var payload *blockchain.Tx
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		blockchain.Mempool().AddPeerTx(payload)
	case MessageNewPeerNotify:
		var payload string
		utility.HandleError(json.Unmarshal(m.Playload, &payload))
		parts := strings.Split(payload, ":")
		AddPeer(parts[0], parts[1], parts[2], false)
	}
}
