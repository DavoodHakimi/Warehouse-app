import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const faNumber = new Intl.NumberFormat('fa-IR')

export function formatNumber(value: number | string | null | undefined): string {
  if (value === null || value === undefined || value === '') return '—'
  const num = typeof value === 'string' ? Number(value) : value
  if (Number.isNaN(num)) return '—'
  return faNumber.format(num)
}

export function formatPrice(
  value: number | string | null | undefined,
  currency = 'ریال',
): string {
  const formatted = formatNumber(value)
  if (formatted === '—') return formatted
  return `${formatted} ${currency}`
}
