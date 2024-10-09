import { log } from '../../utils/log.mjs'
import net from 'net'
import fs from 'node:fs/promises'
import { BaseProcessMessageService } from './base.mjs'

// 这是一个使用 IPC 通道的进程间通信服务类，用于处理来自 Go 进程的消息。
const {
  SOCKET_PATH
} = process.env

/**
 * 使用 IPC 通道的进程间通信服务类，用于处理来自 Go 进程的消息。
 * @class
 * @extends BaseProcessMessageService
 */
export default class IpcProcessMessage extends BaseProcessMessageService {
  #server
  #socket

  /**
   * 构造函数，初始化 IPC 服务。
   */
  constructor () {
    super()
    this.#server = net.createServer((socket) => {
      this.#socket = socket

      log('Client connected')

      socket.on('data', (data) => {
        this.onData(data)
      })

      socket.on('error', (err) => {
        console.error('Socket error:', err)
      })
    })
  }

  /**
   * 初始化 IPC 服务。
   * @returns {Promise<IpcProcessMessage>} 返回初始化后的 IPC 服务实例。
   */
  async init () {
    this.socketPath = SOCKET_PATH
    try {
      if (SOCKET_PATH.startsWith('@')) {
        this.socketPath = `\0${this.socketPath.slice(1)}`
      } else {
        await fs.access(SOCKET_PATH, fs.constants.F_OK)
        await fs.unlink(SOCKET_PATH)
      }
    } catch (err) {
      if (err.code !== 'ENOENT') {
        throw err
      }
    }
    this.#server.listen(this.socketPath, () => {
      log(`Server listening on ${SOCKET_PATH}`)
    })
    return this
  }

  /**
   * 发送数据到 IPC 通道。
   * @param {Buffer} buf 要发送的数据。
   */
  send (buf) {
    if (!this.#socket) {
      log('NO_SOCKET', 'Socket not available for sending data') // 记录更多的上下文信息
      return
    }
    this.#socket.write(buf)
  }

  /**
   * 关闭 IPC 服务并删除监听文件。
   * @returns {Promise<void>}
   */
  async close () {
    this.#server.close()
    // 删除监听文件
    if (SOCKET_PATH.startsWith('/')) {
      try {
        await fs.promises.unlink(this.socketPath)
    } catch (err) {
        console.error('Error deleting socket file:', err) // 添加错误处理
    }
  }
}
