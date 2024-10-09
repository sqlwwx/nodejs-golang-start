package processMessage

import (
	pb "app/messages"
)

type ProcessMessageService interface {
	SendMessage(msg *pb.ProcessMessage)
	Start(onMsg func(*pb.ProcessMessage))
}

func Get(name string) ProcessMessageService {
	switch name {
	case "ipc":
		return &IpcProcessMessageService{}
	case "stdout":
		return &StdoutProcessMessageService{}
	default:
		return &StdoutProcessMessageService{}
	}
}
