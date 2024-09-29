package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/proto"

	pb "app/messages"
	service "app/service"
)

// 全局变量
var (
	// mock,rocket-mq
	MESSAGE_SERVICE_IMPL = flag.String("message-service-impl", "mock", "消息服务实现方式")
	messageService       service.MessageService
)

func initMessageService() {
	flag.Parse()
	messageService, _ = service.GetImplementation(*MESSAGE_SERVICE_IMPL)

	if messageService == nil {
		log.Println("No implementation found ", *MESSAGE_SERVICE_IMPL)
		return
	}
	messageService.Init()
}

func main() {
	log.SetOutput(os.Stderr)
	initMessageService()
	log.Println("Starting go server", *MESSAGE_SERVICE_IMPL)

	reader := bufio.NewReader(os.Stdin)

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
		log.Println("exec startSubscription")
		result := startSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_STOP:
		log.Println("exec stopSubscription")
		result := stopSubscription()
		sendResult(msg.RequestId, result)
	case pb.ProcessMessage_ROCKETMQ_MESSAGE_ACK:
		handleAck(msg.Info)
	}
}

func startSubscription() *pb.ProcessMessage_Info {
	return messageService.Subscribe(func(msg *pb.ProcessMessage) {
		sendMessage(msg)
	})
}

func stopSubscription() *pb.ProcessMessage_Info {
	return messageService.Unsubscribe()
}

func handleAck(info *pb.ProcessMessage_Info) {
	messageService.AckMsg(info)
	if info.GetCode() == 0 {
		log.Println("Message", info.MessageId, "acked successfully", messageService.GetPendingMessageCount())
	} else {
		log.Println("Message", info.MessageId, "failed to ack")
	}
}

func sendMessage(msg *pb.ProcessMessage) {
	data, _ := proto.Marshal(msg)

	header := []byte{0x14, 0x06, 0, 0, 0, 0, 0x15, 0x06}
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(data)))

	os.Stdout.Write(header)
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
