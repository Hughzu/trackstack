const toHex = (buffer: ArrayBuffer) => {
  return Array.from(new Uint8Array(buffer), (byte) => byte.toString(16).padStart(2, '0')).join('')
}

const isPayloadlessMethod = (method: string) => {
  const normalized = method.toUpperCase()
  return normalized === 'GET' || normalized === 'HEAD'
}

const toArrayBuffer = async (body: BodyInit | null | undefined) => {
  if (body == null) {
    return new ArrayBuffer(0)
  }

  if (typeof body === 'string') {
    return new TextEncoder().encode(body).buffer
  }

  if (body instanceof URLSearchParams) {
    return new TextEncoder().encode(body.toString()).buffer
  }

  return new Response(body).arrayBuffer()
}

export const attachPayloadHashHeader = async (headers: Headers, method: string, body?: BodyInit | null) => {
  if (headers.has('x-amz-content-sha256') || isPayloadlessMethod(method)) {
    return
  }

  const payload = await toArrayBuffer(body)
  const digest = await crypto.subtle.digest('SHA-256', payload)

  headers.set('x-amz-content-sha256', toHex(digest))
}
