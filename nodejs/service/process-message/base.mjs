import { log } from '../../utils/log.mjs'
import { ProcessMessage } from '../../utils/proto.mjs'

export class BaseProcessMessageService {
  #onMsg
  #buf = Buffer.alloc(0)
  #offset = 0

  start (onMsg) {
    this.#onMsg = onMsg
  }

  onData (data) {
    this.#buf = Buffer.concat([this.#buf, data])
    if (!this.#onMsg) {
      return
    }
    // 循环查询符合先后缀的数据
    while (this.#offset + 9 <= this.#buf.length) {
      if (this.#buf[this.#offset] !== 20) {
        this.#offset++
        continue
      }
      if (this.#buf[this.#offset + 1] !== 6) {
        this.#offset++
        continue
      }
      if (this.#buf[this.#offset + 6] !== 21) {
        this.#offset++
        continue
      }
      if (this.#buf[this.#offset + 7] !== 6) {
        this.#offset++
        continue
      }

      const length = this.#buf.subarray(this.#offset + 2, this.#offset + 6).readUInt32LE(0)

      if (this.#offset + 8 + length > this.#buf.length) {
        // 数据还未完全到达，等待更多数据
        break
      }

      try {
        const message = ProcessMessage.decode(
          this.#buf.subarray(this.#offset + 8, this.#offset + 8 + length)
        )
        this.#buf = this.#buf.subarray(this.#offset + 8 + length)
        this.#offset = 0
        this.#onMsg(message)
        continue
      } catch (error) {
        console.error(error)
        this.#offset += 2 // 跳过当前解析失败的部分
      }
    }

    // 如果 this.#offset 达到一定大小，或者解析完成后，重置 this.#buf 和 this.#offset
    if (this.#offset >= 1024 || this.#offset >= this.#buf.length) {
      this.#buf = this.#buf.subarray(this.#offset)
      this.#offset = 0
    }
  }

  sendMessage (message) {
    const buffer = ProcessMessage.encode(message).finish()
    log('Sending message:', buffer.length, JSON.stringify(message))
    const lengthBuffer = Buffer.alloc(4)
    lengthBuffer.writeUInt32LE(buffer.length, 0)
    this.send(Buffer.concat([lengthBuffer, buffer]))
  }

  async init () {
    return this
  }
}
