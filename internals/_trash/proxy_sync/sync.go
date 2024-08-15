package proxy

import (
	"fmt"

	"github.com/moov-io/iso8583"
	"github.com/sirupsen/logrus"
)

type (
	ProxyMessageSynchronizer interface {
		AddSync(message *iso8583.Message, channel chan *iso8583.Message)
		Sync(rawMessage []byte) error
	}

	proxyMessageSynchronizer struct {
		iso8583Spec  *iso8583.MessageSpec
		syncChannels map[string]chan *iso8583.Message
	}
)

func NewProxyMessageSynchronizer(iso8583Spec *iso8583.MessageSpec) *proxyMessageSynchronizer {
	return &proxyMessageSynchronizer{
		iso8583Spec:  iso8583Spec,
		syncChannels: make(map[string]chan *iso8583.Message),
	}
}

func (p *proxyMessageSynchronizer) AddSync(message *iso8583.Message, channel chan *iso8583.Message) {
	key := p.syncKey(message)

	p.syncChannels[key] = channel
}

func (p *proxyMessageSynchronizer) Sync(rawMessage []byte) error {
	message := iso8583.NewMessage(p.iso8583Spec)
	if err := message.Unpack(rawMessage); err != nil {
		logrus.WithError(err).Error("failed to open the msg iso")
		return err
	}

	key := p.syncKey(message)

	syncChan, ok := p.syncChannels[key]
	if !ok {
		logrus.WithField("key", key).Warn("sync key was not founded for the message received")
		return fmt.Errorf("sync key not rounded")
	}

	syncChan <- message
	delete(p.syncChannels, key)

	return nil
}

func (p *proxyMessageSynchronizer) syncKey(message *iso8583.Message) string {
	trxDateTime, _ := message.GetField(7).String()
	stun, _ := message.GetField(11).String()
	localTrxTime, _ := message.GetField(12).String()

	return fmt.Sprintf("%v:%v:%v", trxDateTime, stun, localTrxTime)
}
