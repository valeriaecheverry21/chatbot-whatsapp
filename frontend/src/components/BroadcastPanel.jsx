import { useState } from 'react'
import { useWhatsAppApi } from '../hooks/useApi'

export default function BroadcastPanel() {
  const api = useWhatsAppApi()
  const [templateName, setTemplateName] = useState('')
  const [params, setParams] = useState('')
  const [phones, setPhones] = useState('')
  const [result, setResult] = useState(null)
  const [sending, setSending] = useState(false)

  async function handleBroadcast(e) {
    e.preventDefault()
    if (!templateName.trim()) return
    setSending(true)
    setResult(null)

    try {
      const phoneList = phones ? phones.split('\n').map(s => s.trim()).filter(Boolean) : []
      const res = await api.broadcast(templateName, 'es', params ? params.split(',').map(s => s.trim()) : [], phoneList)
      setResult(res)
    } catch (err) {
      setResult({ error: err.message })
    }
    setSending(false)
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <div>
        <h2 className="text-xl font-bold">Broadcast</h2>
        <p className="text-sm text-gray-400">Envía una plantilla a múltiples clientes</p>
      </div>

      <form onSubmit={handleBroadcast} className="space-y-4">
        <div>
          <label className="block text-sm text-gray-400 mb-1">Nombre de plantilla</label>
          <input value={templateName} onChange={e => setTemplateName(e.target.value)}
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600" />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Parámetros (coma)</label>
          <input value={params} onChange={e => setParams(e.target.value)} placeholder="ej. Juan, Promo 2026"
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600" />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">
            Números específicos (uno por línea, o dejar vacío para todos los clientes con opt-in)
          </label>
          <textarea value={phones} onChange={e => setPhones(e.target.value)} rows={4}
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600 resize-none" />
        </div>
        <button disabled={sending}
          className="w-full bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 rounded-xl py-3 font-semibold transition-colors">
          {sending ? 'Enviando...' : 'Enviar broadcast'}
        </button>
      </form>

      {result && (
        <div className={`${result.error ? 'bg-red-900/30 border-red-800' : 'bg-emerald-900/30 border-emerald-800'} border rounded-xl p-4`}>
          {result.error ? (
            <p className="text-red-400 font-semibold">❌ {result.error}</p>
          ) : (
            <>
              <p className="text-emerald-400 font-semibold">✅ Broadcast encolado</p>
              <p className="text-sm text-gray-400 mt-1">
                {result.messages_queued} de {result.total_customers} mensajes encolados
              </p>
            </>
          )}
        </div>
      )}
    </div>
  )
}
