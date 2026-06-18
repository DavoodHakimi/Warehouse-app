// نقش‌های کاربری
export const ROLES: Record<number, string> = {
  1: 'مدیر',
  2: 'سرپرست انبار',
  3: 'انباردار',
  4: 'انباردار - ورود',
  5: 'انباردار - خروج',
}

export const ROLE_OPTIONS = Object.entries(ROLES).map(([id, label]) => ({
  value: Number(id),
  label,
}))

export const CURRENCIES: Record<number, string> = {
  1: 'ریال',
  2: 'دلار',
  3: 'یورو',
}

export const CURRENCY_OPTIONS = Object.entries(CURRENCIES).map(([id, label]) => ({
  value: Number(id),
  label,
}))

export const PARTNER_TYPES: Record<number, string> = {
  1: 'تأمین‌کننده',
  2: 'خریدار',
  3: 'خریدار / فروشنده',
}

export const PARTNER_TYPE_OPTIONS = Object.entries(PARTNER_TYPES).map(
  ([id, label]) => ({ value: Number(id), label }),
)

export const PARTNER_TYPE_LABELS: Record<string, string> = {
  Supplier: 'تأمین‌کننده',
  Customer: 'خریدار',
  Both: 'خریدار / فروشنده',
}

export const ORDER_TYPES: Record<string, string> = {
  sale: 'فروش',
  purchase: 'خرید',
}

export const ORDER_STATUSES: Record<string, string> = {
  Pending: 'در انتظار تأیید',
  Approved: 'تأیید شده',
  Packed: 'بسته‌بندی شده',
  Shipped: 'ارسال شده',
  Waiting: 'در انتظار دریافت',
  Received: 'دریافت شده',
  Canceled: 'لغو شده',
}

export const STATUS_VARIANTS: Record<
  string,
  'default' | 'secondary' | 'outline' | 'destructive' | 'success' | 'warning'
> = {
  Pending: 'warning',
  Approved: 'default',
  Packed: 'secondary',
  Shipped: 'success',
  Waiting: 'warning',
  Received: 'success',
  Canceled: 'destructive',
}
