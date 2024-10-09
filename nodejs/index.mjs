import 'dotenv/config'
import { randomUUID } from 'crypto'
import process from 'node:process'

import { initProcessMessageService } from './service/process-message/index.mjs'
import { log, ProcessMessage } from './utils/index.mjs'
import * as golangProcessUtil from './utils/golang-process.mjs'

// 启动 Go 进程
const golangProcess = await golangProcessUtil.start()

// 创建请求映射
const requestMap = new Map()

// 订阅状态
let subscriptionActive = false
let stopPending = false

const processMessageService = await initProcessMessageService(golangProcess)

processMessageService.start((msg) => {
  handleMessage(msg)
})

let pending = 0

async function handleMessage (message) {
  log('处理消息', message.info)
  let requestId
  let resolve
  switch (message.type) {
    case ProcessMessage.Type.ROCKETMQ_MESSAGE:
    // 处理RocketMQ消息
      log('Received RocketMQ message', pending += 1, message.info.messageId)
      setTimeout(() => {
      // 发送ACK消息
        pending -= 1
        return processMessageService.sendMessage({
          type: ProcessMessage.Type.ROCKETMQ_MESSAGE_ACK,
          info: {
            messageId: message.info.messageId
          }
        })
      }, Math.floor(Math.random() * 30_000))
      break
    case ProcessMessage.Type.RESULT:
    // 处理结果消息
      requestId = message.requestId
      resolve = requestMap.get(requestId)
      if (resolve) {
        resolve(message.info)
        requestMap.delete(requestId)
      }
      break
    default:
      log('未知消息类型', message)
  }
}

/**
 * 发送命令给 Go 进程
 * @param {Object} command 命令对象
 * @returns {Promise} 命令执行结果
 */
async function sendCommand (command) {
  return new Promise((resolve) => {
    // 生成一个唯一的请求 ID
    const requestId = randomUUID()
    command.requestId = requestId
    // 将请求 ID 和解析函数存储在 Map 中
    requestMap.set(requestId, resolve)
    // 发送命令给 Go 进程
    processMessageService.sendMessage(command)
  })
}

/**
 * 启动订阅
 * @returns {Promise} 订阅结果
 */
async function startSubscription () {
  if (subscriptionActive) {
    console.error('Subscription already active')
    return
  }
  subscriptionActive = true
  const result = await sendCommand({
    type: ProcessMessage.Type.START
  })
  log('Start subscription result:', result)
}

/**
 * 停止订阅
 * @returns {Promise} 停止结果
 */
async function stopSubscription () {
  if (!subscriptionActive) {
    console.error('Subscription not active')
    return
  }
  subscriptionActive = false
  const result = await sendCommand({
    type: ProcessMessage.Type.STOP
  })
  log('Stop subscription result:', result)
}

setTimeout(() => {
  // 启动订阅
  startSubscription()
}, 3000)

/**
 * 等待所有消息处理完成
 * @returns {Promise}
 */
const waitDone = async () => {
  return new Promise((resolve) => {
    const timer = setInterval(() => {
      if (pending === 0) {
        clearInterval(timer)
        log('doneMessages')
        resolve()
      }
    }, 500)
  })
}

const handleExit = async (signal) => {
  log('on', signal)
  if (stopPending) {
    return
  }
  stopPending = true
  // 设置 5 分钟超时
  const timer = setTimeout(() => {
    process.exit(0)
  }, 300_000)
  try {
    await stopSubscription()
    await waitDone()
  } catch (error) {
    console.error(error)
  } finally {
    golangProcess.kill('SIGINT')
    clearTimeout(timer)
    console.log(golangProcess.killed)
    processMessageService.close()
    process.exit(0)
  }
}

[
  'SIGHUP',
  'SIGINT',
  'SIGTERM'
].forEach((signal) => {
  process.on(signal, handleExit)
})
