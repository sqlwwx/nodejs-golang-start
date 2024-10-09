package processMessage

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"os"

	pb "app/messages"

	"google.golang.org/protobuf/proto"
)

// 标准输出进程消息服务
type StdoutProcessMessageService struct {
	ProcessMessageService
}

func (s *StdoutProcessMessageService) SendMessage(msg *pb.ProcessMessage) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}

	header := []byte{0x14, 0x06, 0, 0, 0, 0, 0x15, 0x06}
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(data)))

	buf := bytes.NewBuffer(header)
	buf.Write(data)
	_, err = os.Stdout.Write(buf.Bytes())
	if err != nil {
		log.Println("Error writing message to stdout:", err)
	}
}

// 实现进程消息服务接口
func (s *StdoutProcessMessageService) Start(onMsg func(*pb.ProcessMessage)) {
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
		onMsg(&msg)
	}
}
