import { decodeJwt } from "jose"

export function isAccessTokenExpired(token: string): boolean {
  try {
    const payload = decodeJwt(token)
    if (!payload.exp) return true
    // 提前 30 秒认为过期，留出刷新时间
    return Date.now() / 1000 > payload.exp - 30
  } catch {
    return true
  }
}
