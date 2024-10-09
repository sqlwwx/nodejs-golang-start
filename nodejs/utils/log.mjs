export const log = (...args) => console.log('NODEJS', process.pid, new Date().toLocaleString(), ...args)
