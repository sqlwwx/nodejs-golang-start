package message

import (
	"sync"

	pb "app/messages"
)

// 消息结构体
type RocketMQMessage struct {
	Message   string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	MessageId string `protobuf:"bytes,2,opt,name=messageId,proto3" json:"messageId,omitempty"`
}

// MessageService 接口
type MessageService interface {
	WaitDoneAbleInterface
	// 初始化
	Init()
	// 订阅消息
	Subscribe(callback MessageCallback) *pb.ProcessMessage_Info
	// 取消订阅
	Unsubscribe() *pb.ProcessMessage_Info
	AckMsg(info *pb.ProcessMessage_Info)
}

type MessageCallback func(msg *pb.ProcessMessage)

// 获取实现
func Get(name string) MessageService {
	switch name {
	case "rocket-mq":
		return &RocketMQMessageService{
			subscriptionActive: false,
		}
	case "mock":
	default:
	}
	return &MockMessageService{
		subscriptionActive: false,
		messageQueue:       make(chan *RocketMQMessage, 100),
		wg:                 sync.WaitGroup{},
		wgAck:              sync.WaitGroup{},
	}
}
