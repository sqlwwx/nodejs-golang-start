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

type IpcProcessMessageService struct {
	ProcessMessageService,
	conn net.Conn
}

func dialFillZero(network, name string) (net.Conn, error) {
	raw := []byte(name)
	addr_b := make([]byte, 108)
	for i, c := range raw {
		addr_b[i] = c
	}
	addr := string(addr_b)
	return net.Dial(network, addr)
}

func (s *IpcProcessMessageService) getConnection() net.Conn {
	if s.conn == nil {
		socketPath := os.Getenv("SOCKET_PATH")

		if strings.HasPrefix(socketPath, "@") {
			conn, err := dialFillZero("unix", socketPath)
			if err != nil {
				log.Fatal(err)
			} else {
				s.conn = conn
			}

		} else {
			addr, _ := net.ResolveUnixAddr("unix", socketPath)
			conn, err := net.Dial("unix", addr.String())
			if err != nil {
				log.Fatal(err)
			} else {
				s.conn = conn
			}
		}
	}
	return s.conn
}

func (s *IpcProcessMessageService) SendMessage(msg *pb.ProcessMessage) {
	data, _ := proto.Marshal(msg)
	header := []byte{0x14, 0x06, 0, 0, 0, 0, 0x15, 0x06}
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(data)))
	s.getConnection()
	buf := bytes.NewBuffer(header)
	buf.Write(data)
	s.conn.Write(buf.Bytes())
}

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
