import { NavLink } from 'react-router-dom'
import { MessageSquare, Users, LayoutTemplate, Download, Radio, Activity } from 'lucide-react'

const nav = [
  { to: '/', icon: MessageSquare, label: 'Chat' },
  { to: '/plantillas', icon: LayoutTemplate, label: 'Plantillas' },
  { to: '/broadcast', icon: Download, label: 'Broadcast' },
  { to: '/clientes', icon: Users, label: 'Clientes' },
  { to: '/monitor', icon: Radio, label: 'Monitor' },
  { to: '/estado', icon: Activity, label: 'Estado' },
]

export default function Layout({ children }) {
  return (
    <div className="flex h-screen">
      <nav className="hidden md:flex flex-col w-60 bg-gray-900 border-r border-gray-800 p-4">
        <div className="flex items-center gap-3 mb-8 px-3">
          <div className="w-9 h-9 rounded-lg bg-emerald-600 flex items-center justify-center font-bold text-lg">W</div>
          <div>
            <div className="font-bold">WhatsApp Bot</div>
            <div className="text-xs text-gray-400">Admin Panel</div>
          </div>
        </div>
        <div className="flex-1 space-y-1">
          {nav.map(item => (
            <NavLink key={item.to} to={item.to} end={item.to === '/'}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm transition-colors ${
                  isActive ? 'bg-emerald-600/20 text-emerald-400' : 'text-gray-400 hover:bg-gray-800 hover:text-gray-200'
                }`
              }>
              <item.icon size={18} />
              {item.label}
            </NavLink>
          ))}
        </div>
        <div className="text-xs text-gray-600 px-3 pt-4 border-t border-gray-800">
          v1.0.0 • Meta Cloud API
        </div>
      </nav>

      <nav className="md:hidden fixed bottom-0 left-0 right-0 bg-gray-900 border-t border-gray-800 flex z-50">
        {nav.map(item => (
          <NavLink key={item.to} to={item.to} end={item.to === '/'}
            className={({ isActive }) =>
              `flex-1 flex flex-col items-center py-2 text-[10px] transition-colors ${
                isActive ? 'text-emerald-400' : 'text-gray-500'
              }`
            }>
            <item.icon size={18} />
            {item.label}
          </NavLink>
        ))}
      </nav>

      <main className="flex-1 overflow-y-auto p-4 md:p-6 pb-20 md:pb-6">
        {children}
      </main>
    </div>
  )
}
