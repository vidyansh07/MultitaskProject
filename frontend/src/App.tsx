import { Routes, Route } from 'react-router-dom'
import { useAuthStore } from './store/authStore'
import { useEffect } from 'react'

// Layout Components
import Layout from './components/layout/Layout'
import AuthLayout from './components/layout/AuthLayout'
import ProtectedRoute from './components/auth/ProtectedRoute'

// Page Components
import HomePage from './pages/HomePage'
import LoginPage from './pages/auth/LoginPage'
import RegisterPage from './pages/auth/RegisterPage'
import ForgotPasswordPage from './pages/auth/ForgotPasswordPage'
import ResetPasswordPage from './pages/auth/ResetPasswordPage'
import DashboardPage from './pages/DashboardPage'
import ProfilePage from './pages/profile/ProfilePage'
import ChatPage from './pages/chat/ChatPage'
import PostsPage from './pages/posts/PostsPage'
import MarketplacePage from './pages/marketplace/MarketplacePage'
import AIToolsPage from './pages/ai/AIToolsPage'
import SettingsPage from './pages/settings/SettingsPage'
import NotFoundPage from './pages/NotFoundPage'

// Hooks
import { useAuth } from './hooks/useAuth'

function App() {
  const { initialize } = useAuthStore()
  const { isLoading } = useAuth()

  useEffect(() => {
    initialize()
  }, [initialize])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="flex flex-col items-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
          <p className="text-gray-600 font-medium">Loading Multitask Platform...</p>
        </div>
      </div>
    )
  }

  return (
    <Routes>
      {/* Public Routes */}
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
      </Route>

      {/* Authentication Routes */}
      <Route path="/auth" element={<AuthLayout />}>
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="forgot-password" element={<ForgotPasswordPage />} />
        <Route path="reset-password" element={<ResetPasswordPage />} />
      </Route>

      {/* Protected Routes */}
      <Route path="/app" element={<ProtectedRoute><Layout /></ProtectedRoute>}>
        <Route index element={<DashboardPage />} />
        <Route path="profile" element={<ProfilePage />} />
        <Route path="chat" element={<ChatPage />} />
        <Route path="chat/:roomId" element={<ChatPage />} />
        <Route path="posts" element={<PostsPage />} />
        <Route path="marketplace" element={<MarketplacePage />} />
        <Route path="ai-tools" element={<AIToolsPage />} />
        <Route path="settings" element={<SettingsPage />} />
      </Route>

      {/* 404 Route */}
      <Route path="*" element={<NotFoundPage />} />
    </Routes>
  )
}

export default App