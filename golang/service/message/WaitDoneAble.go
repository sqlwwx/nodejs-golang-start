package message

import (
	"log"
	"time"
	"sync"
)

const (
	// 等待完成消息的超时时间
	WAIT_DONE_MSG_TIME = 10 * time.Minute
	// 清理等待完成消息的时间间隔
	CLEAN_WAIT_DONE_MSG_TIME = 5 * time.Minute
)

// WaitDoneAbleInterface 定义了等待完成消息的相关操作接口
type WaitDoneAbleInterface interface {
	startCleanTimeoutMessages()
	addWaitDoneMessage(id string)
	DoneMessage(id string)
	GetPendingMessageCount() int
}

// WaitDoneAble 结构体用于管理等待完成的消息
type WaitDoneAble struct {
	waitDoneMessages map[string]time.Time
	mu               sync.RWMutex // 添加读写锁
}

// startCleanTimeoutMessages 启动一个定时器，定期清理超时的等待完成消息
func (w *WaitDoneAble) startCleanTimeoutMessages() {
	w.waitDoneMessages = make(map[string]time.Time)
	ticker := time.NewTicker(CLEAN_WAIT_DONE_MSG_TIME)
	go func() {
		for range ticker.C {
			log.Println("Cleaning wait done messages", w.waitDoneMessages)
			now := time.Now()
			w.mu.Lock() // 加写锁
			for id, timestamp := range w.waitDoneMessages {
				if now.Sub(timestamp) > WAIT_DONE_MSG_TIME {
					delete(w.waitDoneMessages, id)
				}
			}
			w.mu.Unlock() // 释放写锁
		}
	}()
}

// addWaitDoneMessage 添加一个等待完成的消息
func (w *WaitDoneAble) addWaitDoneMessage(id string) {
	w.mu.Lock() // 加写锁
	w.waitDoneMessages[id] = time.Now()
	w.mu.Unlock() // 释放写锁
}

// DoneMessage 标记一个消息为已完成
func (w *WaitDoneAble) DoneMessage(id string) {
	w.mu.Lock() // 加写锁
	delete(w.waitDoneMessages, id)
	w.mu.Unlock() // 释放写锁
}

// GetPendingMessageCount 获取当前等待完成的消息数量
func (w *WaitDoneAble) GetPendingMessageCount() int {
	w.mu.RLock() // 加读锁
	count := len(w.waitDoneMessages)
	w.mu.RUnlock() // 释放读锁
	return count
}
