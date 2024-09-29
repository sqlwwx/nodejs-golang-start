package service

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

const (
	Topic         = "NORMAL"
	ConsumerGroup = "LOCAL_DEV"
	Endpoint      = "ep-bp1ie0935b53ed066d22.epsrv-bp1g2k84c7kid2igyiun.cn-hangzhou.privatelink.aliyuncs.com:8080"
)

var (
	AccessKey               = os.Getenv("ROCKETMQ_ACCESS_KEY")
	SecretKey               = os.Getenv("ROCKETMQ_SECRET_KEY")
	awaitDuration           = time.Second * 100
	maxMessageNum     int32 = 16
	invisibleDuration       = time.Second * 400
)

type RocketMQMessageServiceImpl struct {
	WaitDoneAble
	subscriptionActive bool
	consumer           rmq_client.SimpleConsumer
	pendingMsgMap      map[string]*rmq_client.MessageView
}

func (m *RocketMQMessageServiceImpl) Init() {
	log.SetOutput(os.Stderr)
	// os.Setenv("mq.consoleAppender.enabled", "true")
	// rmq_client.ResetLogger()
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
	}
	err = consumer.Start()
	if err != nil {
		log.Println("start consumer error:", err)
	}
	m.consumer = consumer
	m.startCleanTimeoutMessages()
}

func (m *RocketMQMessageServiceImpl) Subscribe(callback MessageCallback) *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		return &pb.ProcessMessage_Info{Code: 1, Message: "Subscription already active"}
	}
	m.subscriptionActive = true
	go func() {
		for m.subscriptionActive {
			mvs, err := m.consumer.Receive(context.TODO(), maxMessageNum, invisibleDuration)
			if err != nil && !strings.Contains(err.Error(), "MESSAGE_NOT_FOUND") {
				log.Println("loadMsg error:", err)
			}
			for _, mv := range mvs {
				MessageId := mv.GetMessageId()
				Message := string(mv.GetBody())
				log.Println("onMsg", MessageId, Message)
				rocketMQMessage := &pb.ProcessMessage{
					RequestId: uuid.New().String(),
					Type:      pb.ProcessMessage_ROCKETMQ_MESSAGE,
					Info: &pb.ProcessMessage_Info{
						MessageId: MessageId,
						Message:   Message,
					},
				}
				callback(rocketMQMessage)
				m.pendingMsgMap[mv.GetMessageId()] = mv
				m.addWaitDoneMessage(mv.GetMessageId())
			}
			time.Sleep(time.Second * 1)
		}
	}()
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription started"}
}

func (m *RocketMQMessageServiceImpl) Unsubscribe() *pb.ProcessMessage_Info {
	if m.subscriptionActive {
		m.subscriptionActive = false
		log.Println("start consumer gracefulStop")
		m.consumer.GracefulStop()
		log.Println("done consumer gracefulStop")
	}
	return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription stopped"}
}

func (m *RocketMQMessageServiceImpl) AckMsg(info *pb.ProcessMessage_Info) {
	log.Println("ack", info)
	m.consumer.Ack(context.TODO(), m.pendingMsgMap[info.MessageId])
	m.DoneMessage(info.MessageId)
}
