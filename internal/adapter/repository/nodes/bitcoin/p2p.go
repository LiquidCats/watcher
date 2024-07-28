package bitcoin

import (
	"context"
)

type P2pClient struct {
}

func (receiver *P2pClient) Subscribe(ctx context.Context) error {
	//verack := make(chan struct{})
	//peerCfg := &peer.Config{
	//	ChainParams:     &chaincfg.SimNetParams,
	//	Services:        0,
	//	TrickleInterval: time.Second * 10,
	//	Listeners: peer.MessageListeners{
	//		OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
	//			verack <- struct{}{}
	//		},
	//		OnMemPool: func(p *peer.Peer, msg *wire.MsgMemPool) {
	//			fmt.Println("outbound: received mempool")
	//		},
	//		OnTx: func(p *peer.Peer, msg *wire.MsgTx) {
	//			fmt.Println("outbound: received tx")
	//		},
	//		OnBlock: func(p *peer.Peer, msg *wire.MsgBlock, buf []byte) {
	//			fmt.Println("outbound: received block")
	//		},
	//	},
	//	AllowSelfConns: true,
	//}

	return nil
}

func (receiver *P2pClient) Unsubscribe(ctx context.Context) error {
	return nil
}
