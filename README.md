# 🤖 WhatsApp Chatbot - Business Messaging Platform

A powerful, enterprise-grade WhatsApp messaging platform built with **Go backend** and a **responsive React frontend**. Send messages, manage templates, broadcast campaigns, and monitor webhooks in real-time.

---

## 📋 Table of Contents

- [Features](#features)
- [Tech Stack](#tech-stack)
- [Frontend - Web App](#frontend---web-app-)
- [Backend - API](#backend---api-)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Documentation](#documentation)

---

## ✨ Features

### 🎯 Core Capabilities
- **Real-time Messaging** - Send and receive WhatsApp messages instantly
- **Message Templates** - Pre-approved templates for faster messaging
- **Broadcast Campaigns** - Send bulk messages to customers with opt-in management
- **Live Webhooks** - Real-time event monitoring with SSE (Server-Sent Events)
- **Customer Management** - Track and manage customer communication preferences
- **Health Monitoring** - API status dashboard for backend health checks

---

## 🛠️ Tech Stack

| Component | Technology | Version |
|-----------|-----------|---------|
| **Backend** | Go | 1.21+ |
| **Frontend** | React | 18+ |
| **Build Tool** | Vite | 5+ |
| **Styling** | Tailwind CSS | 3+ |
| **HTTP Client** | Axios | Latest |
| **Responsive Design** | Mobile-first CSS Grid | — |

---

## 🎨 Frontend - Web App

### 📱 Responsive & Cross-Platform

A **modern, fully responsive web application** that works seamlessly on **mobile, tablet, and desktop** without any native app installation. Built with cutting-edge web technologies for maximum performance and user experience.

```
✅ Mobile-First Design     — Optimized for phones & tablets
✅ Zero Installation       — Just open in any modern browser
✅ PWA Ready              — Can work offline (with service workers)
✅ Real-time Updates      — Live webhook monitoring via SSE
✅ Touch-Friendly UI      — Optimized touch targets for mobile
✅ Fast Performance       — Vite hot module replacement during development
```

#### 🎯 Key Features

**📬 Chat Interface**
- WhatsApp-style messaging UI
- Send and receive messages in real-time
- Message history and threading
- Auto-refresh and live notifications

**📋 Template Management**
- Test and preview approved message templates
- Template builder with variables
- Quick template selection for faster messaging

**📢 Broadcast Panel**
- Manage customer segments
- Schedule bulk message campaigns
- Track delivery status and analytics
- Opt-in/opt-out management

**👥 Customer Management**
- View all registered customers
- Filter by opt-in status
- Search and sort functionality
- Customer communication preferences

**🔔 Live Webhook Monitor**
- Real-time event streaming (Server-Sent Events)
- Incoming message notifications
- Delivery status updates
- Activity logs and debugging

**🏥 API Health Dashboard**
- Backend service status
- Response time metrics
- Connection indicators

#### 🎨 UI Components

| Component | Purpose |
|-----------|---------|
| **Layout.jsx** | Main navigation (sidebar on desktop, bottom tabs on mobile) |
| **ChatView.jsx** | WhatsApp-style messaging interface |
| **TemplateTester.jsx** | Preview & test message templates |
| **BroadcastPanel.jsx** | Bulk messaging & campaign management |
| **CustomerList.jsx** | Customer database & opt-in tracker |
| **WebhookMonitor.jsx** | Live SSE event streaming |
| **ApiStatus.jsx** | Backend health indicators |

#### 🚀 Tech Highlights

- **Vite**: Lightning-fast development server with HMR
- **Tailwind CSS**: Utility-first styling for rapid UI development
- **Responsive Grid**: Mobile-first CSS Grid layout
- **Custom Hooks**: `useApi.js` for centralized HTTP client management
- **Real-time SSE**: WebSocket-like streaming without WebSocket overhead

#### 📱 Mobile vs Desktop

**Mobile Layout** 📱
```
┌─────────────────┐
│   Navigation    │  ← Sticky bottom tabs
├─────────────────┤
│                 │
│   Main Content  │  ← Full-width responsive
│                 │
├─────────────────┤
│ 🏠 📋 📢 👥 🔔  │  ← Bottom navigation
└─────────────────┘
```

**Desktop Layout** 🖥️
```
┌──────────┬─────────────────────────┐
│          │                         │
│ Sidebar  │                         │
│          │   Main Content          │
│ 🏠 📋 📢 │   (Full Width)          │
│ 👥 🔔    │                         │
│          │                         │
└──────────┴─────────────────────────┘
```

#### 💻 How to Run Frontend

```bash
# Install dependencies
cd frontend
npm install

# Development server (with hot reload)
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

**Access:** Open `http://localhost:5173` in your browser

#### 🔌 HTTP Client Integration

The `useApi.js` hook abstracts all API communication:

```javascript
const { data, loading, error, request } = useApi();

// Call backend endpoints
const customers = await request('/api/v1/customers');
const sent = await request('/api/v1/messages', 'POST', { ... });
```

#### 🔄 Future: React Native Migration

The modular component architecture allows easy migration to **React Native** for native iOS/Android apps:

```
React Web (Vite)  ──→  Shared Logic & Hooks
                  ──→  React Native (Mobile)
                  ──→  Web (PWA)
```

Simply replace Tailwind styles with React Native Stylesheet and swap web components for RN components.

---

## 🔧 Backend - API

### Go RESTful API

**Key Endpoints:**

```
GET  /api/v1/customers           # List all customers
POST /api/v1/messages            # Send message
GET  /api/v1/messages/:id        # Get message details
POST /api/v1/templates/test      # Test template
POST /api/v1/broadcast/send      # Send bulk campaign
GET  /api/v1/webhooks/stream     # SSE webhook stream
GET  /api/v1/health              # Health check
```

### 📊 Real-time Streaming

- **SSE (Server-Sent Events)** for low-latency webhook delivery
- No WebSocket complexity
- Works through standard HTTP proxies
- Auto-reconnect capability

---

## 🚀 Getting Started

### Prerequisites

- **Node.js 16+** (for frontend)
- **Go 1.21+** (for backend)
- **npm or yarn** (package manager)

### Setup

1. **Clone Repository**
   ```bash
   git clone https://github.com/valeriaecheverry21/chatbot-whatsapp.git
   cd chatbot-whatsapp
   ```

2. **Frontend Setup**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

3. **Backend Setup**
   ```bash
   cd ..
   go mod download
   go run main.go
   ```

4. **Access Application**
   - Frontend: `http://localhost:5173`
   - Backend API: `http://localhost:8080/api/v1`

---

## 📁 Project Structure

```
chatbot-whatsapp/
├── frontend/                       # React + Vite + Tailwind
│   ├── src/
│   │   ├── components/            # React components
│   │   │   ├── ChatView.jsx        # 💬 WhatsApp-style chat
│   │   │   ├── TemplateTester.jsx  # 📋 Template preview
│   │   │   ├── BroadcastPanel.jsx  # 📢 Bulk messaging
│   │   │   ├── CustomerList.jsx    # 👥 Customer database
│   │   │   ├── WebhookMonitor.jsx  # 🔔 Live events
│   │   │   ├── ApiStatus.jsx       # 🏥 Health monitor
│   │   │   └── Layout.jsx          # 🧭 Navigation
│   │   ├── hooks/
│   │   │   └── useApi.js           # 🔌 HTTP client
│   │   ├── App.jsx                 # 🎯 Router
│   │   └── main.jsx
│   ├── vite.config.js              # Build config + proxy
│   ├── tailwind.config.js           # Style config
│   └── package.json
│
├── controllers/                    # Go API handlers
│   ├── sse.go                      # Real-time streaming
│   └── customer_controller.go      # Customer endpoints
│
├── models/                         # Data structures
├── services/                       # Business logic
├── middleware/                     # HTTP middleware
├── config/                         # Configuration
├── main.go                         # Entry point
└── go.mod                          # Go dependencies
```

---

## 📚 Documentation

- **API Docs**: See `docs/api.md`
- **Setup Guide**: See `docs/setup.md`
- **WhatsApp Integration**: See `docs/whatsapp-api.md`

---

## 🔐 Security

- ✅ Environment variables for sensitive config
- ✅ API key validation on all endpoints
- ✅ CORS policies configured
- ✅ Request validation & sanitization
- ✅ Opt-in consent management

---

## 📝 License

This project is licensed under the MIT License. See `LICENSE` file for details.

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📧 Support

For questions or issues, please open a GitHub issue or contact the maintainers.

---

## 🙏 Acknowledgments

Built with ❤️ for seamless WhatsApp business messaging.

**Happy Messaging! 🚀**
