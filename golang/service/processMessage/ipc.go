package processMessage

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"strings"

	pb "app/messages"

	"google.golang.org/protobuf/proto"
)

// IpcProcessMessageService 结构体实现了 ProcessMessageService 接口，并包含一个 net.Conn 连接
type IpcProcessMessageService struct {
	ProcessMessageService
	conn net.Conn
}

// dialFillZero 函数用于填充地址字节数组，确保地址长度为108字节
func dialFillZero(network, name string) (net.Conn, error) {
	raw := []byte(name)
	addr := make([]byte, 108)
	copy(addr, raw)
	return net.Dial(network, string(addr))
}

// getConnection 方法用于获取或创建与 IPC 的连接
func (s *IpcProcessMessageService) getConnection() net.Conn {
	if s.conn == nil {
		socketPath := os.Getenv("SOCKET_PATH")
		var err error

		if strings.HasPrefix(socketPath, "@") {
			s.conn, err = dialFillZero("unix", socketPath)
		} else {
			addr, _ := net.ResolveUnixAddr("unix", socketPath)
			s.conn, err = net.Dial("unix", addr.String())
		}

		if err != nil {
			log.Fatal(err)
		}
	}
	return s.conn
}

// SendMessage 方法用于发送消息到 IPC
func (s *IpcProcessMessageService) SendMessage(msg *pb.ProcessMessage) {
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling message:", err)
		return
	}
	header := []byte{0x14, 0x06, 0, 0, 0, 0, 0x15, 0x06}
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(data)))
	s.getConnection()
	buf := bytes.NewBuffer(header)
	buf.Write(data)
	_, err = s.conn.Write(buf.Bytes())
	if err != nil {
		log.Println("Error writing message:", err)
	}
}

// Start 方法用于启动消息处理循环
func (s *IpcProcessMessageService) Start(onMsg func(*pb.ProcessMessage)) {
	s.getConnection()
	reader := bufio.NewReader(s.conn)
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
