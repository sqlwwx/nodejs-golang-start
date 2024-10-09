import { BaseProcessMessageService } from './base.mjs'

export default class StdioProcessMessage extends BaseProcessMessageService {
  #golangProcess

  constructor (golangProcess) {
    super()
    this.#golangProcess = golangProcess
    // 处理 Go 进程输出
    this.#golangProcess.stdout.on('data', (data) => {
      this.onData(data)
    })
  }

  send (buf) {
    this.#golangProcess.stdin.write(buf)
  }

  async close () {}
}
