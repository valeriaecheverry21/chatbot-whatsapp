import { useState, useEffect, useRef } from 'react'

export default function WebhookMonitor() {
  const [events, setEvents] = useState([])
  const [connected, setConnected] = useState(false)
  const bottomRef = useRef(null)

  useEffect(() => {
    const es = new EventSource('/api/v1/webhooks/stream')
    es.onopen = () => setConnected(true)
    es.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data)
        setEvents(prev => [{ ...data, id: Date.now() }, ...prev].slice(0, 100))
      } catch { /* ignore */ }
    }
    es.onerror = () => setConnected(false)
    return () => es.close()
  }, [])

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [events])

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold">Monitor Webhook</h2>
          <p className="text-sm text-gray-400">Eventos entrantes en tiempo real</p>
        </div>
        <span className={`flex items-center gap-2 text-sm ${connected ? 'text-emerald-400' : 'text-red-400'}`}>
          <span className={`w-2 h-2 rounded-full ${connected ? 'bg-emerald-400 animate-pulse' : 'bg-red-400'}`} />
          {connected ? 'Conectado' : 'Desconectado'}
        </span>
      </div>

      <div className="bg-gray-900 rounded-xl h-[500px] overflow-y-auto p-4 space-y-2 font-mono text-xs">
        {events.length === 0 && (
          <div className="text-center text-gray-500 mt-20">
            <div className="text-3xl mb-2">📡</div>
            <p>Esperando eventos...</p>
            <p className="text-xs mt-1">Los webhooks aparecerán aquí cuando Meta envíe mensajes</p>
            <p className="text-xs text-gray-600 mt-2">¿Servidor backend corriendo?</p>
          </div>
        )}
        {events.map(ev => (
          <div key={ev.id} className="bg-gray-800 rounded-lg p-3 border-l-4 border-emerald-600">
            <div className="flex justify-between text-gray-400 mb-1">
              <span>{ev.type || 'webhook'}</span>
              <span>{new Date(ev.id).toLocaleTimeString()}</span>
            </div>
            <pre className="text-gray-300 whitespace-pre-wrap break-all">
              {JSON.stringify(ev.data || ev, null, 2)}
            </pre>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>
    </div>
  )
}
