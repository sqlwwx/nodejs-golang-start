# nodejs-golang-start

这是一个 Node.js 和 Go 语言混合开发的项目，使用了 protobuf 进行通信。

## 项目结构

- `main.go`: Go 语言的主程序，负责处理消息队列和发送结果。
- `index.mjs`: Node.js 的主程序，负责启动 Go 进程、发送命令和处理结果。
- `messages.proto`: protobuf 文件，定义了消息格式。
- `package.json`: Node.js 项目的配置文件。

## 依赖

- `protobufjs`: 用于解析 protobuf 文件。

## 安装

1. 安装依赖：`pnpm i`
2. 构建 go 程序：`go run main.go`
3. 运行 Node.js 程序：`node index.mjs`

## 生成 protobuf 代码

```

protoc --go_out=./golang proto/messages.proto

```

## 使用

1. 启动 Node.js 程序后，会自动启动 Go 进程。
2. 可以通过 Node.js 程序发送命令给 Go 进程。
3. Go 进程会处理命令并发送结果给 Node.js 程序。
4. Node.js 程序会处理结果并输出到控制台。

## 进程通信方式

- [x] 标准输入输出 简单 数据传输量有限
- [ ] 文件共享/内存文件 适用大数据量传输/持久化 慢/文件锁/并发问题
- [ ] 网络服务 分布式 网络延迟/复杂度高
- [ ] 共享内存

## TOTO

- [ ] 使用共享内存方式实现进程间通信

## License

MIT License
