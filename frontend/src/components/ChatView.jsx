import { useState, useRef, useEffect } from 'react'
import { useWhatsAppApi } from '../hooks/useApi'

const BOT_NUMBER = '521234567890'

const quickReplies = [
  { label: 'Hola', msg: 'Hola, quiero información' },
  { label: 'Horarios', msg: '¿Cuál es el horario de atención?' },
  { label: 'Precios', msg: '¿Me pueden dar los precios?' },
  { label: 'Gracias', msg: 'Gracias, muy amable' },
]

export default function ChatView() {
  const api = useWhatsAppApi()
  const [messages, setMessages] = useState([])
  const [input, setInput] = useState('')
  const [phone, setPhone] = useState(BOT_NUMBER)
  const [sending, setSending] = useState(false)
  const bottomRef = useRef(null)

  useEffect(() => { bottomRef.current?.scrollIntoView({ behavior: 'smooth' }) }, [messages])

  async function handleSend(e) {
    e.preventDefault()
    if (!input.trim() || !phone.trim()) return

    const userMsg = { id: Date.now(), from: 'me', text: input.trim(), time: new Date().toLocaleTimeString() }
    setMessages(prev => [...prev, userMsg])
    setInput('')
    setSending(true)

    try {
      const res = await api.sendText(phone, userMsg.text)
      setMessages(prev => [...prev, {
        id: Date.now() + 1, from: 'bot', text: `✅ Enviado (ID: ${res.whatsapp_id?.slice(0, 15)}...)`,
        time: new Date().toLocaleTimeString(), sent: true,
      }])
    } catch (err) {
      setMessages(prev => [...prev, {
        id: Date.now() + 1, from: 'bot', text: `❌ Error: ${err.message}`,
        time: new Date().toLocaleTimeString(), error: true,
      }])
    }
    setSending(false)
  }

  return (
    <div className="flex flex-col h-[calc(100vh-8rem)] max-w-2xl mx-auto">
      <div className="flex items-center gap-3 p-4 border-b border-gray-800">
        <div className="w-10 h-10 rounded-full bg-emerald-600 flex items-center justify-center text-lg font-bold">
          B
        </div>
        <div>
          <div className="font-semibold">ChatBot WhatsApp</div>
          <div className="text-xs text-gray-400">{phone}</div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4 space-y-3">
        {messages.length === 0 && (
          <div className="text-center text-gray-500 mt-12">
            <div className="text-5xl mb-4">💬</div>
            <p>Envía un mensaje de prueba</p>
            <p className="text-sm mt-2">El chatbot responderá simulando el envío</p>
          </div>
        )}
        {messages.map(msg => (
          <div key={msg.id} className={`flex ${msg.from === 'me' ? 'justify-end' : 'justify-start'}`}>
            <div className={`max-w-[80%] rounded-2xl px-4 py-2 ${
              msg.from === 'me'
                ? 'bg-emerald-700 rounded-br-sm'
                : msg.error
                  ? 'bg-red-900/50 rounded-bl-sm'
                  : 'bg-gray-800 rounded-bl-sm'
            }`}>
              <p className="text-sm">{msg.text}</p>
              <p className={`text-[10px] mt-1 ${msg.from === 'me' ? 'text-emerald-200' : 'text-gray-400'}`}>
                {msg.time}
                {msg.sent && <span className="ml-1">✓✓</span>}
              </p>
            </div>
          </div>
        ))}
        <div ref={bottomRef} />
      </div>

      <div className="p-3 border-t border-gray-800">
        <div className="flex gap-2 mb-2 overflow-x-auto pb-1">
          {quickReplies.map(q => (
            <button key={q.label} onClick={() => setInput(q.msg)}
              className="shrink-0 text-xs bg-gray-800 hover:bg-gray-700 rounded-full px-3 py-1.5 transition-colors">
              {q.label}
            </button>
          ))}
        </div>
        <form onSubmit={handleSend} className="flex gap-2">
          <input
            value={phone} onChange={e => setPhone(e.target.value)}
            placeholder="Número destino"
            className="w-36 bg-gray-800 rounded-xl px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-emerald-600"
          />
          <input
            value={input} onChange={e => setInput(e.target.value)}
            placeholder="Escribe un mensaje..."
            className="flex-1 bg-gray-800 rounded-xl px-4 py-2 text-sm outline-none focus:ring-2 focus:ring-emerald-600"
          />
          <button disabled={sending || !input.trim()}
            className="bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 rounded-xl px-4 py-2 font-medium transition-colors">
            {sending ? '...' : '→'}
          </button>
        </form>
      </div>
    </div>
  )
}
