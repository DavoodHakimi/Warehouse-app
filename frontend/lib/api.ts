export const API_BASE =
  process.env.NEXT_PUBLIC_API_URL?.replace(/\/$/, '') ||
  'http://localhost:8080/api/v1'

const TOKEN_KEY = 'wms_token'

export function getToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
}

export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

type RequestOptions = {
  method?: string
  body?: unknown
  auth?: boolean
}

export async function apiFetch<T = unknown>(
  path: string,
  { method = 'GET', body, auth = true }: RequestOptions = {},
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }

  if (auth) {
    const token = getToken()
    if (token) headers['Authorization'] = `Bearer ${token}`
  }

  let res: Response
  try {
    res = await fetch(`${API_BASE}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    })
  } catch {
    throw new ApiError(
      'اتصال به سرور برقرار نشد. مطمئن شوید بک‌اند روی localhost:8080 در حال اجراست.',
      0,
    )
  }

  const text = await res.text()
  let data: unknown = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
  }

  if (!res.ok) {
    const message =
      (data && typeof data === 'object' && 'error' in data
        ? String((data as Record<string, unknown>).error)
        : data && typeof data === 'object' && 'message' in data
          ? String((data as Record<string, unknown>).message)
          : null) || errorMessageForStatus(res.status)
    throw new ApiError(message, res.status)
  }

  return data as T
}

function errorMessageForStatus(status: number): string {
  switch (status) {
    case 400:
      return 'درخواست نامعتبر است.'
    case 401:
      return 'دسترسی غیرمجاز. لطفاً دوباره وارد شوید.'
    case 403:
      return 'شما مجوز انجام این عملیات را ندارید.'
    case 404:
      return 'مورد درخواستی یافت نشد.'
    case 409:
      return 'تداخل در داده‌ها.'
    case 500:
      return 'خطای داخلی سرور.'
    default:
      return 'خطای ناشناخته رخ داد.'
  }
}
