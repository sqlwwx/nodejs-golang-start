package service

import (
	"sync"
	"time"

	"github.com/google/uuid"

	pb "app/messages"
)

type MockMessageServiceImpl struct {
	WaitDoneAble
	subscriptionActive bool
	messageQueue       chan *RocketMQMessage
	wg                 sync.WaitGroup
	wgAck              sync.WaitGroup
}

func (m *MockMessageServiceImpl) Init() {
	m.startCleanTimeoutMessages()
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if m.subscriptionActive {
				messageId := uuid.New().String()
				m.messageQueue <- &RocketMQMessage{
					Message:   "Simulated RocketMQ Message",
					MessageId: messageId,
				}
			}
		}
	}()
}

func (m *MockMessageServiceImpl) Subscribe(callback MessageCallback) *pb.ProcessMessage_Info {
	if !m.subscriptionActive {
		m.subscriptionActive = true
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			for msg := range m.messageQueue {
				rocketMQMessage := &pb.ProcessMessage{
					RequestId: uuid.New().String(),
					Type:      pb.ProcessMessage_ROCKETMQ_MESSAGE,
					Info: &pb.ProcessMessage_Info{
						MessageId: msg.MessageId,
						Message:   msg.Message,
					},
				}
				m.wgAck.Add(1)
				callback(rocketMQMessage)
				m.addWaitDoneMessage(msg.MessageId)
			}
		}()
		return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription started"}
	}
	return &pb.ProcessMessage_Info{Code: 1, Message: "Subscription already active"}
}

func (m *MockMessageServiceImpl) Unsubscribe() *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		m.subscriptionActive = false
		close(m.messageQueue)
		m.wg.Wait()
	}
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription stopped"}
}

func (m *MockMessageServiceImpl) AckMsg(info *pb.ProcessMessage_Info) {
	m.DoneMessage(info.MessageId)
}
