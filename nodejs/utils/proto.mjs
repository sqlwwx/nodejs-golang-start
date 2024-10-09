import protobuf from 'protobufjs'

export const root = protobuf.loadSync('proto/messages.proto')
export const ProcessMessage = root.lookupType('messages.ProcessMessage')
