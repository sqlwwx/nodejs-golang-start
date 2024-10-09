package message

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	pb "app/messages"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	"github.com/google/uuid"
)

var (
	Topic             = os.Getenv("ROCKETMQ_TOPIC")
	ConsumerGroup     = os.Getenv("ROCKETMQ_CONSUMER_GROUP")
	Endpoint          = os.Getenv("ROCKETMQ_ENDPOINT")
	AccessKey         = os.Getenv("ROCKETMQ_ACCESS_KEY")
	SecretKey         = os.Getenv("ROCKETMQ_SECRET_KEY")
	awaitDuration     = 100 * time.Second
	maxMessageNum int32 = 16
	invisibleDuration = 400 * time.Second
)

type RocketMQMessageService struct {
	WaitDoneAble
	subscriptionActive bool
	consumer           rmq_client.SimpleConsumer
	pendingMsgMap      map[string]*rmq_client.MessageView
}

func (m *RocketMQMessageService) Init() {
	log.SetOutput(os.Stderr)
	m.pendingMsgMap = make(map[string]*rmq_client.MessageView)
	consumer, err := rmq_client.NewSimpleConsumer(&rmq_client.Config{
		Endpoint:      Endpoint,
		ConsumerGroup: ConsumerGroup,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    AccessKey,
			AccessSecret: SecretKey,
		},
	},
		rmq_client.WithAwaitDuration(awaitDuration),
		rmq_client.WithSubscriptionExpressions(map[string]*rmq_client.FilterExpression{
			Topic: rmq_client.SUB_ALL,
		}),
	)
	if err != nil {
		log.Println("init consumer:", err)
		return
	}
	if err := consumer.Start(); err != nil {
		log.Println("start consumer error:", err)
		return
	}
	m.consumer = consumer
	m.startCleanTimeoutMessages()
}

func (m *RocketMQMessageService) Subscribe(callback MessageCallback) *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		return &pb.ProcessMessage_Info{Code: 1, Message: "Subscription already active"}
	}
	m.subscriptionActive = true
	go func() {
		for m.subscriptionActive {
			mvs, err := m.consumer.Receive(context.TODO(), maxMessageNum, invisibleDuration)
			if err != nil && !strings.Contains(err.Error(), "MESSAGE_NOT_FOUND") {
				log.Println("loadMsg error:", err)
				continue
			}
			for _, mv := range mvs {
				messageId := mv.GetMessageId()
				message := string(mv.GetBody())
				log.Println("onMsg", messageId, message)
				rocketMQMessage := &pb.ProcessMessage{
					RequestId: uuid.New().String(),
					Type:      pb.ProcessMessage_ROCKETMQ_MESSAGE,
					Info: &pb.ProcessMessage_Info{
						MessageId: messageId,
						Message:   message,
					},
				}
				callback(rocketMQMessage)
				m.pendingMsgMap[messageId] = mv
				m.addWaitDoneMessage(messageId)
			}
			time.Sleep(time.Second)
		}
	}()
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription started"}
}

func (m *RocketMQMessageService) Unsubscribe() *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		m.subscriptionActive = false
		log.Println("start consumer gracefulStop")
		if err := m.consumer.GracefulStop(); err != nil {
			log.Println("consumer gracefulStop error:", err)
		}
		log.Println("done consumer gracefulStop")
	}
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription stopped"}
}

func (m *RocketMQMessageService) AckMsg(info *pb.ProcessMessage_Info) {
	log.Println("ack", info)
	if err := m.consumer.Ack(context.TODO(), m.pendingMsgMap[info.MessageId]); err != nil {
		log.Printf("ack error: %v, messageId: %s", err, info.MessageId) // 记录更多的上下文信息
	}
	m.DoneMessage(info.MessageId)
}
