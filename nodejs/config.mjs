import * as dotenv from 'dotenv'

/**
 * 加载环境变量配置文件，并将其导出为默认对象。
 * @returns {Object} 包含环境变量的对象
 * @throws {Error} 如果加载环境变量文件失败，则抛出错误
 */
const loadEnvConfig = () => {
  const result = dotenv.config()
  if (result.error) {
    throw result.error
  }
  return { ...result.parsed }
}

export default loadEnvConfig()
