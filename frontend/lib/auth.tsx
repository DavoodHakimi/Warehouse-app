'use client'

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from 'react'
import { apiFetch, clearToken, getToken, setToken } from './api'
import { can, type Permission } from './permissions'

export type AuthUser = {
  userId: number
  username: string
  companyId: number
  role: number
}

type AuthContextValue = {
  user: AuthUser | null
  loading: boolean
  login: (userName: string, password: string) => Promise<void>
  logout: () => void
  can: (permission: Permission) => boolean
}

const AuthContext = createContext<AuthContextValue | null>(null)

function decodeToken(token: string): AuthUser | null {
  try {
    const payload = token.split('.')[1]
    const json = JSON.parse(
      decodeURIComponent(
        atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
          .split('')
          .map((c) => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
          .join(''),
      ),
    )

    const num = (v: unknown) => (typeof v === 'string' ? Number(v) : (v as number))

    if (json.exp && Date.now() / 1000 > Number(json.exp)) {
      return null
    }

    return {
      userId: num(json.user_id),
      username: String(json.username ?? ''),
      companyId: num(json.company_id),
      role: num(json.role),
    }
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = getToken()
    if (token) {
      const decoded = decodeToken(token)
      if (decoded) {
        setUser(decoded)
      } else {
        clearToken()
      }
    }
    setLoading(false)
  }, [])

  const login = useCallback(async (userName: string, password: string) => {
    const data = await apiFetch<{ token: string }>('/auth/login', {
      method: 'POST',
      auth: false,
      body: { user_name: userName, password },
    })
    setToken(data.token)
    const decoded = decodeToken(data.token)
    if (!decoded) {
      throw new Error('توکن دریافتی نامعتبر است.')
    }
    setUser(decoded)
  }, [])

  const logout = useCallback(() => {
    clearToken()
    setUser(null)
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      loading,
      login,
      logout,
      can: (permission: Permission) => can(user?.role, permission),
    }),
    [user, loading, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
