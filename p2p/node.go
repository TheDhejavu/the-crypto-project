package p2p

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

// NodeRoomBufSize is the number of incoming messages to buffer for each topic.
const NodeRoomBufSize = 128

// NodeRoom represents a subsnoderoomiption to a single PubSub topic. Messages
// can be published to the topic with NodeRoom.Publish, and received
// messages are pushed to the Messages channel.
type NodeRoom struct {
	ctx   context.Context
	pub   *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	roomName string
	self     peer.ID
	Data     chan *NodeContent
}
type NodeContent struct {
	Message string
	NodeID  string
	SendTo  string
	Payload []byte
}

// JoinNodeRoom tries to subsnoderoomibe to the PubSub topic for the network, returning
// a NodeRoom on success.
func JoinNodeRoom(ctx context.Context, pub *pubsub.PubSub, selfID peer.ID, roomName string, subscribe bool) (*NodeRoom, error) {
	// join the pubsub topic
	topic, err := pub.Join(topicName(roomName))
	if err != nil {
		return nil, err
	}

	// and subsnoderoomibe to it
	var sub *pubsub.Subscription

	if subscribe {
		sub, err = topic.Subscribe()
		if err != nil {
			return nil, err
		}
	} else {
		sub = nil
	}

	noderoom := &NodeRoom{
		ctx:      ctx,
		pub:      pub,
		topic:    topic,
		sub:      sub,
		self:     selfID,
		roomName: roomName,
		Data:     make(chan *NodeContent, NodeRoomBufSize),
	}

	go noderoom.readLoop()
	return noderoom, nil
}

func (node *NodeRoom) ListPeers() []peer.ID {
	return node.pub.ListPeers(topicName(node.roomName))
}

func topicName(roomName string) string {
	return "node-room:" + roomName
}

// Publish sends a message to the pubsub topic.
func (node *NodeRoom) Publish(message string, sendTo string) error {
	m := NodeContent{
		Message: message,
		NodeID:  shortID(node.self),
		SendTo:  sendTo,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return node.topic.Publish(node.ctx, msgBytes)
}

func (node *NodeRoom) readLoop() {
	if node.sub == nil {
		return
	}
	for {
		content, err := node.sub.Next(node.ctx)
		if err != nil {
			close(node.Data)
			return
		}
		// only forward messages delivered by others
		if content.ReceivedFrom == node.self {
			continue
		}

		nd := new(NodeContent)
		err = json.Unmarshal(content.Data, nd)
		if err != nil {
			continue
		}

		if nd.SendTo != "" && nd.SendTo != shortID(node.self) {
			continue
		}
		// send valid messages onto the Messages channel
		// fmt.Println(nd)
		node.Data <- nd
	}
}
