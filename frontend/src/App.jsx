import { Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import ChatView from './components/ChatView'
import TemplateTester from './components/TemplateTester'
import BroadcastPanel from './components/BroadcastPanel'
import CustomerList from './components/CustomerList'
import WebhookMonitor from './components/WebhookMonitor'
import ApiStatus from './components/ApiStatus'

export default function App() {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<ChatView />} />
        <Route path="/plantillas" element={<TemplateTester />} />
        <Route path="/broadcast" element={<BroadcastPanel />} />
        <Route path="/clientes" element={<CustomerList />} />
        <Route path="/monitor" element={<WebhookMonitor />} />
        <Route path="/estado" element={<ApiStatus />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </Layout>
  )
}
