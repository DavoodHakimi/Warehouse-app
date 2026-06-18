'use client'

import { useCallback, useEffect, useState } from 'react'
import { apiFetch, ApiError } from './api'

type State<T> = {
  data: T | null
  error: string | null
  loading: boolean
  refetch: () => void
}

export function useApi<T>(path: string | null): State<T> {
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState<boolean>(!!path)
  const [tick, setTick] = useState(0)

  const refetch = useCallback(() => setTick((t) => t + 1), [])

  useEffect(() => {
    if (!path) {
      setLoading(false)
      return
    }
    let active = true
    setLoading(true)
    setError(null)
    apiFetch<T>(path)
      .then((res) => {
        if (active) setData(res)
      })
      .catch((err) => {
        if (active)
          setError(err instanceof ApiError ? err.message : 'خطا در دریافت داده.')
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => {
      active = false
    }
  }, [path, tick])

  return { data, error, loading, refetch }
}
