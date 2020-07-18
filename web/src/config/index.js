export const BACKEND_PORT = process.env.PORT || 4000
export const WEBSOCKET_SCHEME = process.env.NODE_ENV === 'production' ? 'wss' : 'ws'
