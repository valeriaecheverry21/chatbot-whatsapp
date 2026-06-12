const BASE = ''

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json', ...options.headers },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || `HTTP ${res.status}`)
  }
  return res.json()
}

export function useWhatsAppApi() {
  return {
    health: () => request('/health'),

    sendText: (phone, message, preview = false) =>
      request('/api/v1/messages/text', {
        method: 'POST',
        body: JSON.stringify({ phone, message, preview }),
      }),

    sendImage: (phone, url, caption = '') =>
      request('/api/v1/messages/image', {
        method: 'POST',
        body: JSON.stringify({ phone, url, caption }),
      }),

    sendTemplate: (phone, name, language = 'es', params = []) =>
      request('/api/v1/messages/template', {
        method: 'POST',
        body: JSON.stringify({ phone, name, language, params }),
      }),

    queueMessage: (phone, message) =>
      request('/api/v1/messages/queue', {
        method: 'POST',
        body: JSON.stringify({ phone, message }),
      }),

    broadcast: (templateName, language = 'es', params = [], phones = []) =>
      request('/api/v1/broadcast', {
        method: 'POST',
        body: JSON.stringify({ template_name: templateName, language, params, phones }),
      }),
  }
}
