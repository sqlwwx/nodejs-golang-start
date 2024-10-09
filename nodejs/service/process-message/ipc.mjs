import { log } from '../../utils/log.mjs'
import net from 'net'
import fs from 'node:fs/promises'
import { BaseProcessMessageService } from './base.mjs'

const {
  SOCKET_PATH
} = process.env

export default class IpcProcessMessage extends BaseProcessMessageService {
  #server
  #socket

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

  send (buf) {
    if (!this.#socket) {
      log('NO_SOCKET')
      return
    }
    this.#socket.write(buf)
  }

  async close () {
    this.#server.close()
    // 删除监听文件
    if (SOCKET_PATH.startsWith('/')) {
      return fs.promises.unlink(this.socketPath)
    }
  }
}
