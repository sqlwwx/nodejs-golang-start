import { spawn } from 'child_process'
import { log } from './log.mjs'

const { NODE_ENV } = process.env

export const start = async () => {
  // 根据环境变量设置命令行参数
  const commandArgs = NODE_ENV === 'production' ? ['./app', []] : ['go', ['run', 'main.go']]

  log('start golang-process', commandArgs)

  // 启动 Go 进程
  const golangProcess = spawn(...commandArgs, {
    detached: true,
    cwd: 'golang',
    env: {
      ...process.env,
      RPC_PID: process.pid
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
  return golangProcess
}
