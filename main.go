package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	pb "nodejs-golang-start/messages"
)

// 全局变量
var (
	subscriptionActive = false
	messageQueue       = make(chan *RocketMQMessage, 100)
	wg                 sync.WaitGroup
	wgAck              sync.WaitGroup
	waitDoneMessages   = make(map[string]time.Time)
)

// 常量定义
const (
	// WAIT_DONE_MSG_TIME       = 20 * time.Minute
	// CLEAN_WAIT_DONE_MSG_TIME = 10 * time.Minute
	WAIT_DONE_MSG_TIME       = 10 * time.Minute
	CLEAN_WAIT_DONE_MSG_TIME = 5 * time.Minute
)

// RocketMQ 消息结构体
type RocketMQMessage struct {
	Message   string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	MessageId string `protobuf:"bytes,2,opt,name=messageId,proto3" json:"messageId,omitempty"`
}

func cleanWaitDoneMessages() {
	log.Println("Cleaning wait done messages", waitDoneMessages)
	now := time.Now()
	for id, timestamp := range waitDoneMessages {
		if now.Sub(timestamp) > WAIT_DONE_MSG_TIME {
			delete(waitDoneMessages, id)
		}
	}
}

func main() {
	// 设置日志输出
	log.SetOutput(os.Stderr)

	log.Println("Starting go server")

	reader := bufio.NewReader(os.Stdin)

	// 启动定时器, 清理超时消息
	ticker := time.NewTicker(CLEAN_WAIT_DONE_MSG_TIME)
	go func() {
		for range ticker.C {
			cleanWaitDoneMessages()
		}
	}()

	// 启动消息队列
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if subscriptionActive {
				messageId := uuid.New().String()
				messageQueue <- &RocketMQMessage{
					Message:   "Simulated RocketMQ Message",
					MessageId: messageId,
				}
			}
		}
	}()

	// 处理消息
	for {
		// 读取消息长度
		var length uint32
		err := binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			log.Println("Error reading message length:", err)
			continue
		}

		// 读取消息
		buffer := make([]byte, length)

		_, err = io.ReadFull(reader, buffer)

		if err != nil {
			log.Println("Error reading message:", err)
			continue
		}

		// 解析 protobuf 消息
		var msg pb.ProcessMessage
		err = proto.Unmarshal(buffer, &msg)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		// 处理消息
		handleMessage(&msg)
	}
}

func handleMessage(msg *pb.ProcessMessage) {
	switch msg.Type {
	case pb.ProcessMessage_START:
		log.Println("START")
		result := startSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_STOP:
		result := stopSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_ROCKETMQ_MESSAGE_ACK:
		handleAck(msg.Info)
	}
}

func startSubscription() *pb.ProcessMessage_Info {
	if !subscriptionActive {
		subscriptionActive = true
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range messageQueue {
				rocketMQMessage := &pb.ProcessMessage{
					RequestId: uuid.New().String(),
					Type:      pb.ProcessMessage_ROCKETMQ_MESSAGE,
					Info: &pb.ProcessMessage_Info{
						MessageId: msg.MessageId,
						Message:   msg.Message,
					},
				}
				wgAck.Add(1)
				sendMessage(rocketMQMessage)
				waitDoneMessages[msg.MessageId] = time.Now()
			}
		}()
		return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription started"}
	}
	return &pb.ProcessMessage_Info{Code: 1, Message: "Subscription already active"}
}

func stopSubscription() *pb.ProcessMessage_Info {
	if subscriptionActive {
		subscriptionActive = false
		close(messageQueue)
		wg.Wait()
		return &pb.ProcessMessage_Info{Code: 0, Message: "Subscription stopped "}
	}
	return &pb.ProcessMessage_Info{Code: 1, Message: "Subscription not active"}
}

func handleAck(info *pb.ProcessMessage_Info) {
	delete(waitDoneMessages, info.MessageId)
	if info.GetCode() == 0 {
		log.Println("Message", info.MessageId, "acked successfully", len(waitDoneMessages))
	} else {
		log.Println("Message", info.MessageId, "failed to ack")
	}
}

func sendMessage(msg *pb.ProcessMessage) {
	data, _ := proto.Marshal(msg)
	length := uint32(len(data))
	binary.Write(os.Stdout, binary.LittleEndian, length)
	os.Stdout.Write(data)
}

func sendResult(requestId string, info *pb.ProcessMessage_Info) {
	msg := &pb.ProcessMessage{
		RequestId: requestId,
		Type:      pb.ProcessMessage_RESULT,
		Info:      info,
	}
	sendMessage(msg)
}
