import { useState } from 'react'
import { useWhatsAppApi } from '../hooks/useApi'

export default function TemplateTester() {
  const api = useWhatsAppApi()
  const [phone, setPhone] = useState('521234567890')
  const [name, setName] = useState('')
  const [params, setParams] = useState('')
  const [result, setResult] = useState(null)
  const [error, setError] = useState('')
  const [sending, setSending] = useState(false)

  async function handleSend(e) {
    e.preventDefault()
    if (!phone.trim() || !name.trim()) return
    setSending(true)
    setResult(null)
    setError('')

    try {
      const res = await api.sendTemplate(phone, name, 'es', params ? params.split(',').map(s => s.trim()) : [])
      setResult(res)
    } catch (err) {
      setError(err.message)
    }
    setSending(false)
  }

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <div>
        <h2 className="text-xl font-bold">Probar Plantilla</h2>
        <p className="text-sm text-gray-400">Envía una plantilla aprobada de Meta a un número</p>
      </div>

      <form onSubmit={handleSend} className="space-y-4">
        <div>
          <label className="block text-sm text-gray-400 mb-1">Número destino</label>
          <input value={phone} onChange={e => setPhone(e.target.value)}
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600" />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Nombre de la plantilla</label>
          <input value={name} onChange={e => setName(e.target.value)} placeholder="ej. bienvenida_oferta"
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600" />
        </div>
        <div>
          <label className="block text-sm text-gray-400 mb-1">Parámetros (separados por coma)</label>
          <input value={params} onChange={e => setParams(e.target.value)} placeholder="ej. Juan, Oferta 50%, 30 días"
            className="w-full bg-gray-800 rounded-xl px-4 py-2.5 outline-none focus:ring-2 focus:ring-emerald-600" />
        </div>
        <button disabled={sending}
          className="w-full bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 rounded-xl py-3 font-semibold transition-colors">
          {sending ? 'Enviando...' : 'Enviar plantilla'}
        </button>
      </form>

      {result && (
        <div className="bg-emerald-900/30 border border-emerald-800 rounded-xl p-4">
          <p className="text-emerald-400 font-semibold">✅ Enviada</p>
          <p className="text-sm text-gray-400 mt-1">WhatsApp ID: {result.whatsapp_id}</p>
        </div>
      )}
      {error && (
        <div className="bg-red-900/30 border border-red-800 rounded-xl p-4">
          <p className="text-red-400 font-semibold">❌ Error</p>
          <p className="text-sm text-gray-400 mt-1">{error}</p>
        </div>
      )}
    </div>
  )
}
