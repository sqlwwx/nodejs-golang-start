import { spawn } from 'child_process'
import { randomUUID } from 'crypto'
import { Buffer } from 'buffer'
import process from 'node:process'
import protobuf from 'protobufjs'

const { NODE_ENV } = process.env

// 加载 protobuf 文件
const root = protobuf.loadSync('messages.proto')
const ProcessMessage = root.lookupType('messages.ProcessMessage')

// 根据环境变量设置命令行参数
const commandArgs = NODE_ENV === 'production' ? ['./nodejs-golang-start', []] : ['go', ['run', 'main.go']]

// 启动 Go 进程
const golangProcess = spawn(...commandArgs, {
  detached: true
})

// 创建请求映射
const requestMap = new Map()

// 订阅状态
let subscriptionActive = false
let stopPending = false

const log = (...args) => console.log(process.pid, ...args)

// 处理 Go 进程输出
golangProcess.stdout.on('data', (data) => {
  // 解析 protobuf 消息
  let offset = 0
  while (offset < data.length) {
    const length = data.readUInt32LE(offset)
    const buffer = data.slice(offset + 4, offset + 4 + length)
    const message = ProcessMessage.decode(buffer)
    handleMessage(message)
    offset += 4 + length
  }
})

// 处理 Go 进程错误输出
golangProcess.stderr.on('data', (data) => {
  console.error('go', golangProcess.pid, data.toString())
})

// 处理 Go 进程退出
golangProcess.on('exit', (code) => {
  console.log(`Go 进程退出，退出码：${code}`)
})

// 处理 Go 进程错误
golangProcess.on('error', (err) => {
  console.error(`Go 进程错误：${err}`)
})

let pending = 0

async function handleMessage (message) {
  console.log(process.pid, 'handleMessage', JSON.stringify(message))
  if (message.type === ProcessMessage.Type.ROCKETMQ_MESSAGE) {
    log('Received RocketMQ message', pending += 1, message.info.messageId)
    setTimeout(() => {
      pending -= 1
      return sendMessage({
        type: ProcessMessage.Type.ROCKETMQ_MESSAGE_ACK,
        info: {
          messageId: message.info.messageId
        }
      })
    }, Math.floor(Math.random() * 3 * 60_000))
  }
  if (message.type === ProcessMessage.Type.RESULT) {
    const requestId = message.requestId
    const resolve = requestMap.get(requestId)
    if (resolve) {
      resolve(message.info)
      requestMap.delete(requestId)
    }
  }
}

function sendMessage (message) {
  const buffer = ProcessMessage.encode(message).finish()
  log('Sending message:', buffer.length, JSON.stringify(message))
  const lengthBuffer = Buffer.alloc(4)
  lengthBuffer.writeUInt32LE(buffer.length, 0)
  golangProcess.stdin.write(lengthBuffer)
  golangProcess.stdin.write(buffer)
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
    sendMessage(command)
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

// 启动订阅
startSubscription()

/**
 * 等待所有消息处理完成
 * @returns {Promise}
 */
const waitDone = async () => {
  return new Promise((resolve) => {
    const timer = setInterval(() => {
      if (pending === 0) {
        clearInterval(timer)
        console.log('doneMessages')
        resolve()
      }
    }, 500)
  })
}

const handleExit = async (signal) => {
  console.log('on', signal)
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
