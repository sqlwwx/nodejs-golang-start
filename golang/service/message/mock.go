package message

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	pb "app/messages"
)

// 模拟消息服务实现
type MockMessageService struct {
	WaitDoneAble
	subscriptionActive bool
	messageQueue       chan *RocketMQMessage
	wg                 sync.WaitGroup
	wgAck              sync.WaitGroup
	count              int
}

func (m *MockMessageService) Init() {
	m.startCleanTimeoutMessages()
	m.count = 0
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if m.subscriptionActive && m.count < 10 {
				m.count += 1
				messageId := uuid.New().String()
				m.messageQueue <- &RocketMQMessage{
					Message:   fmt.Sprintf("message %d", m.count),
					MessageId: messageId,
				}
			}
		}
	}()
}

func (m *MockMessageService) Subscribe(callback MessageCallback) *pb.ProcessMessage_Info {
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

func (m *MockMessageService) Unsubscribe() *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		m.subscriptionActive = false
		close(m.messageQueue)
		m.wg.Wait()
	}
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription stopped"}
}

func (m *MockMessageService) AckMsg(info *pb.ProcessMessage_Info) {
	m.DoneMessage(info.MessageId)
}
