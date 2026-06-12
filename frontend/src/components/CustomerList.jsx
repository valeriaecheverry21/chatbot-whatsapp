import { useState, useEffect } from 'react'

export default function CustomerList() {
  const [customers, setCustomers] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/customers')
      .then(res => res.json())
      .then(data => setCustomers(data.customers || []))
      .catch(() => setCustomers([]))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div>
        <h2 className="text-xl font-bold">Clientes</h2>
        <p className="text-sm text-gray-400">Clientes registrados con opt-in</p>
      </div>

      {loading ? (
        <div className="text-center py-12 text-gray-500">Cargando...</div>
      ) : customers.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <div className="text-4xl mb-3">👥</div>
          <p>No hay clientes registrados</p>
          <p className="text-sm mt-1">Los clientes aparecen cuando reciben mensajes vía webhook</p>
          <p className="text-xs text-gray-600 mt-4">¿PostgreSQL conectado?</p>
        </div>
      ) : (
        <div className="space-y-2">
          {customers.map(c => (
            <div key={c.id} className="bg-gray-900 rounded-xl p-4 flex items-center justify-between">
              <div>
                <p className="font-medium">{c.name || 'Sin nombre'}</p>
                <p className="text-sm text-gray-400">{c.phone}</p>
                <p className="text-xs text-gray-500">{c.created_at}</p>
              </div>
              <span className={`px-3 py-1 rounded-full text-xs font-medium ${c.opt_in ? 'bg-emerald-900 text-emerald-300' : 'bg-gray-800 text-gray-400'}`}>
                {c.opt_in ? 'Opt-in ✓' : 'No opt-in'}
              </span>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
