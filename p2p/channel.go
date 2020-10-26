package p2p

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p-core/peer"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

const ChannelBufSize = 128

type Channel struct {
	ctx   context.Context
	pub   *pubsub.PubSub
	topic *pubsub.Topic
	sub   *pubsub.Subscription

	channelName string
	self        peer.ID
	Content     chan *ChannelContent
}

type ChannelContent struct {
	Message  string
	SendFrom string
	SendTo   string
	Payload  []byte
}

func JoinChannel(ctx context.Context, pub *pubsub.PubSub, selfID peer.ID, channelName string, subscribe bool) (*Channel, error) {

	topic, err := pub.Join(topicName(channelName))
	if err != nil {
		return nil, err
	}

	
	var sub *pubsub.Subscription

	if subscribe {
		sub, err = topic.Subscribe()
		if err != nil {
			return nil, err
		}
	} else {
		sub = nil
	}

	Channel := &Channel{
		ctx:         ctx,
		pub:         pub,
		topic:       topic,
		sub:         sub,
		self:        selfID,
		channelName: channelName,
		Content:     make(chan *ChannelContent, ChannelBufSize),
	}

	go Channel.readLoop()

	return Channel, nil
}

func (ch *Channel) ListPeers() []peer.ID {
	return ch.pub.ListPeers(topicName(ch.channelName))
}

func (channel *Channel) Publish(message string, payload []byte, SendTo string) error {
	m := ChannelContent{
		Message:  message,
		SendFrom: ShortID(channel.self),
		SendTo:   SendTo,
		Payload:  payload,
	}
	msgBytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return channel.topic.Publish(channel.ctx, msgBytes)
}

func (channel *Channel) readLoop() {
	if channel.sub == nil {
		return
	}
	for {
		content, err := channel.sub.Next(channel.ctx)
		if err != nil {
			close(channel.Content)
			return
		}
		// only forward messages delivered by others
		if content.ReceivedFrom == channel.self {
			continue
		}

		NewContent := new(ChannelContent)
		err = json.Unmarshal(content.Data, NewContent)
		if err != nil {
			continue
		}

		if NewContent.SendTo != "" && NewContent.SendTo != channel.self.Pretty() {
			continue
		}

		// send valid messages onto the Messages channel
		channel.Content <- NewContent
	}
}

func topicName(channelName string) string {
	return "channel:" + channelName
}
