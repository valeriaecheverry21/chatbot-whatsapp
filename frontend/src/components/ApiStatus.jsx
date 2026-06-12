import { useState, useEffect } from 'react'

export default function ApiStatus() {
  const [status, setStatus] = useState(null)
  const [error, setError] = useState(null)

  useEffect(() => {
    function check() {
      fetch('/health')
        .then(res => res.json())
        .then(setStatus)
        .catch(err => setError(err.message))
    }
    check()
    const interval = setInterval(check, 10000)
    return () => clearInterval(interval)
  }, [])

  return (
    <div className="max-w-lg mx-auto space-y-6">
      <div>
        <h2 className="text-xl font-bold">Estado del Servicio</h2>
        <p className="text-sm text-gray-400">Backend chatbot-whatsapp</p>
      </div>

      {error ? (
        <div className="bg-red-900/30 border border-red-800 rounded-xl p-6 text-center">
          <div className="text-4xl mb-3">🔴</div>
          <p className="text-red-400 font-semibold">Backend no responde</p>
          <p className="text-sm text-gray-400 mt-2">{error}</p>
          <p className="text-xs text-gray-500 mt-4">Ejecutá `go run .` en la carpeta del backend</p>
        </div>
      ) : !status ? (
        <div className="bg-gray-900 rounded-xl p-6 text-center">
          <div className="animate-spin text-4xl mb-3">⏳</div>
          <p className="text-gray-400">Verificando conexión...</p>
        </div>
      ) : (
        <div className="space-y-4">
          <div className="bg-emerald-900/30 border border-emerald-800 rounded-xl p-4 text-center">
            <div className="text-4xl mb-2">🟢</div>
            <p className="text-emerald-400 font-semibold">Servicio activo</p>
            <p className="text-sm text-gray-400">{status.service}</p>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div className={`rounded-xl p-4 ${status.checks?.database ? 'bg-emerald-900/30 border border-emerald-800' : 'bg-red-900/30 border border-red-800'}`}>
              <p className="text-sm text-gray-400">Base de datos</p>
              <p className={`text-lg font-bold ${status.checks?.database ? 'text-emerald-400' : 'text-red-400'}`}>
                {status.checks?.database ? '✓' : '✗'}
              </p>
            </div>
            <div className={`rounded-xl p-4 ${status.checks?.redis ? 'bg-emerald-900/30 border border-emerald-800' : 'bg-red-900/30 border border-red-800'}`}>
              <p className="text-sm text-gray-400">Redis</p>
              <p className={`text-lg font-bold ${status.checks?.redis ? 'text-emerald-400' : 'text-red-400'}`}>
                {status.checks?.redis ? '✓' : '✗'}
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
