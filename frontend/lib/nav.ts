import type { Permission } from './permissions'

export type NavItem = {
  href: string
  label: string
  icon: string
  permission: Permission
}

export const NAV_ITEMS: NavItem[] = [
  { href: '/dashboard', label: 'داشبورد', icon: 'LayoutDashboard', permission: 'products.read' },
  { href: '/dashboard/orders', label: 'سفارش‌ها', icon: 'ClipboardList', permission: 'orders.read' },
  { href: '/dashboard/products', label: 'محصولات', icon: 'Package', permission: 'products.read' },
  { href: '/dashboard/partners', label: 'شرکای تجاری', icon: 'Handshake', permission: 'partners.read' },
  { href: '/dashboard/users', label: 'کاربران', icon: 'Users', permission: 'users.read' },
]
