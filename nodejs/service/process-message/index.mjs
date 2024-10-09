const {
  PROCESS_MESSAGE = 'stdout'
} = process.env

// 初始化消息处理服务
export const initProcessMessageService = async (golangProcess) => {
  const { default: ProcessMessageService } = await import(`./${PROCESS_MESSAGE}.mjs`)
  return new ProcessMessageService(golangProcess).init()
}
