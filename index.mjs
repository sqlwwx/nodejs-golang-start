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
const commandArgs = NODE_ENV === 'production' ? ['./nodejs-golang-start', []] : ['go', ['run', './main.go']]
// const commandArgs = ['go', ['run', './test/mq.go']]

commandArgs[1].push(...process.argv.slice(2))

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

let goBuf = Buffer.alloc(0)
let offset = 0

// 处理 Go 进程输出
golangProcess.stdout.on('data',
  (data) => {
    goBuf = Buffer.concat([goBuf, data])
    // 循环查询符合先后缀的数据
    while (offset + 9 <= goBuf.length) {
      if (goBuf[offset] !== 20) {
        offset++
        continue
      }
      if (goBuf[offset + 1] !== 6) {
        offset++
        continue
      }
      if (goBuf[offset + 6] !== 21) {
        offset++
        continue
      }
      if (goBuf[offset + 7] !== 6) {
        offset++
        continue
      }

      const length = goBuf.subarray(offset + 2, offset + 6).readUInt32LE(0)

      if (offset + 8 + length > goBuf.length) {
        // 数据还未完全到达，等待更多数据
        break
      }

      try {
        const message = ProcessMessage.decode(
          goBuf.subarray(offset + 8, offset + 8 + length)
        )
        goBuf = goBuf.subarray(offset + 8 + length)
        offset = 0
        handleMessage(message)
        continue
      } catch (error) {
        console.error(error)
        offset += 2 // 跳过当前解析失败的部分
      }
    }

    // 如果 offset 达到一定大小，或者解析完成后，重置 goBuf 和 offset
    if (offset >= 1024 || offset >= goBuf.length) {
      goBuf = goBuf.subarray(offset)
      offset = 0
    }
  })

// 处理 Go 进程错误输出
golangProcess.stderr.on('data', (data) => {
  data.toString().split('\n').filter(s => s)
    .map(s => console.log(golangProcess.pid, s))
})

// 处理 Go 进程退出
golangProcess.on('exit', (code) => {
  log(`Go 进程退出，退出码：${code}`)
})

// 处理 Go 进程错误
golangProcess.on('error', (err) => {
  console.error(`Go 进程错误：${err}`)
})

let pending = 0

async function handleMessage (message) {
  log('handleMessage', message.info)
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
