import { BaseProcessMessageService } from './base.mjs'

/**
 * StdioProcessMessage 类继承自 BaseProcessMessageService，用于处理与 Go 进程的输入输出交互。
 */
export default class StdioProcessMessage extends BaseProcessMessageService {
  /**
   * @private
   * @type {import('child_process').ChildProcess}
   */
  #golangProcess

  /**
   * 构造函数，初始化 StdioProcessMessage 实例。
   * @param {import('child_process').ChildProcess} golangProcess - Go 进程实例
   */
  constructor (golangProcess) {
    super()
    this.#golangProcess = golangProcess
    // 处理 Go 进程输出
    this.#golangProcess.stdout.on('data', (data) => {
      this.onData(data)
    })
  }

  /**
   * 向 Go 进程发送数据。
   * @param {Buffer} buf - 要发送的数据
   */
  send (buf) {
    this.#golangProcess.stdin.write(buf)
  }

  /**
   * 关闭与 Go 进程的连接。
   * @returns {Promise<void>}
   */
  async close () {}
}
