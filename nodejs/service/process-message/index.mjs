const {
  PROCESS_MESSAGE = 'stdout'
} = process.env

export const initProcessMessageService = async (golangProcess) => {
  const { default: ProcessMessageService } = await import(`./${PROCESS_MESSAGE}.mjs`)
  return new ProcessMessageService(golangProcess).init()
}
