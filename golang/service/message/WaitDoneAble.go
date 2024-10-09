package message

import (
	"log"
	"time"
)

const (
	WAIT_DONE_MSG_TIME       = 10 * time.Minute
	CLEAN_WAIT_DONE_MSG_TIME = 5 * time.Minute
)

type WaitDoneAbleInterface interface {
	startCleanTimeoutMessages()
	addWaitDoneMessage(id string)
	DoneMessage(id string)
	GetPendingMessageCount() int
}

type WaitDoneAble struct {
	waitDoneMessages map[string]time.Time
}

func (w *WaitDoneAble) startCleanTimeoutMessages() {
	w.waitDoneMessages = make(map[string]time.Time)
	ticker := time.NewTicker(CLEAN_WAIT_DONE_MSG_TIME)
	go func() {
		for range ticker.C {
			log.Println("Cleaning wait done messages", w.waitDoneMessages)
			now := time.Now()
			for id, timestamp := range w.waitDoneMessages {
				if now.Sub(timestamp) > WAIT_DONE_MSG_TIME {
					delete(w.waitDoneMessages, id)
				}
			}
		}
	}()
}

func (w *WaitDoneAble) addWaitDoneMessage(id string) {
	w.waitDoneMessages[id] = time.Now()
}

func (w *WaitDoneAble) DoneMessage(id string) {
	delete(w.waitDoneMessages, id)
}

func (w *WaitDoneAble) GetPendingMessageCount() int {
	return len(w.waitDoneMessages)
}
