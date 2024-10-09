package main

import (
	"log"
	"os"

	pb "app/messages"
	messageServicePkg "app/service/message"
	"app/service/processMessage"
)

// 全局变量
var (
	// 消息服务类型，可选值：mock, rocket-mq
	MESSAGE_SERVICE_TYPE string
	// 处理消息类型，可选值：stdout, ipc
	PROCESS_MESSAGE_TYPE  string
	messageService        messageServicePkg.MessageService
	processMessageService processMessage.ProcessMessageService
)

func init() {
	log.SetOutput(os.Stderr)
	MESSAGE_SERVICE_TYPE = os.Getenv("MESSAGE_SERVICE")
	PROCESS_MESSAGE_TYPE = os.Getenv("PROCESS_MESSAGE")
	if MESSAGE_SERVICE_TYPE == "" {
		MESSAGE_SERVICE_TYPE = "mock"
	}
	if PROCESS_MESSAGE_TYPE == "" {
		PROCESS_MESSAGE_TYPE = "stdout"
	}
}

func main() {
	messageService = messageServicePkg.Get(MESSAGE_SERVICE_TYPE)
	messageService.Init()

	log.Println("Starting go server", MESSAGE_SERVICE_TYPE, PROCESS_MESSAGE_TYPE, os.Getenv("RPC_PID"))

	processMessageService = processMessage.Get(PROCESS_MESSAGE_TYPE)

	processMessageService.Start(func(msg *pb.ProcessMessage) {
		handleProcessMessage(msg)
	})
}

func handleProcessMessage(msg *pb.ProcessMessage) {
	switch msg.Type {
	case pb.ProcessMessage_START:
		log.Println("执行启动订阅")
		result := startSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_STOP:
		log.Println("执行取消订阅")
		result := stopSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_ROCKETMQ_MESSAGE_ACK:
		handleAck(msg.Info)
	default:
		log.Printf("未知的消息类型: %v", msg.Type)
	}
}

// 启动订阅
func startSubscription() *pb.ProcessMessage_Info {
	return messageService.Subscribe(func(msg *pb.ProcessMessage) {
		processMessageService.SendMessage(msg)
	})
}

// 取消订阅
func stopSubscription() *pb.ProcessMessage_Info {
	return messageService.Unsubscribe()
}

// 处理确认消息
func handleAck(info *pb.ProcessMessage_Info) {
	messageService.AckMsg(info)
	if info.GetCode() == 0 {
		log.Println("消息", info.MessageId, "确认成功", messageService.GetPendingMessageCount())
	} else {
		log.Println("消息", info.MessageId, "确认失败")
	}
}

func sendResult(requestId string, info *pb.ProcessMessage_Info) {
	msg := &pb.ProcessMessage{
		RequestId: requestId,
		Type:      pb.ProcessMessage_RESULT,
		Info:      info,
	}
	processMessageService.SendMessage(msg)
}
