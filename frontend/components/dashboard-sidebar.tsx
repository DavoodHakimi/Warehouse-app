'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import {
  LayoutDashboard,
  ClipboardList,
  Package,
  Handshake,
  Users,
  Warehouse,
  LogOut,
  type LucideIcon,
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAuth } from '@/lib/auth'
import { NAV_ITEMS } from '@/lib/nav'
import { ROLES } from '@/lib/constants'
import { Button } from '@/components/ui/button'

const ICONS: Record<string, LucideIcon> = {
  LayoutDashboard,
  ClipboardList,
  Package,
  Handshake,
  Users,
}

export function DashboardSidebar({ onNavigate }: { onNavigate?: () => void }) {
  const pathname = usePathname()
  const { user, logout, can } = useAuth()

  const items = NAV_ITEMS.filter((item) => can(item.permission))

  return (
    <div className="flex h-full flex-col bg-sidebar">
      <div className="flex items-center gap-3 border-b border-sidebar-border px-5 py-5">
        <div className="flex size-10 items-center justify-center rounded-xl bg-sidebar-primary text-sidebar-primary-foreground">
          <Warehouse className="size-5" />
        </div>
        <div className="leading-tight">
          <p className="text-sm font-bold text-sidebar-foreground">مدیریت انبار</p>
          <p className="text-xs text-muted-foreground">سامانه جامع موجودی</p>
        </div>
      </div>

      <nav className="flex flex-1 flex-col gap-1 overflow-y-auto p-3">
        {items.map((item) => {
          const Icon = ICONS[item.icon] ?? Package
          const active =
            item.href === '/dashboard'
              ? pathname === '/dashboard'
              : pathname.startsWith(item.href)
          return (
            <Link
              key={item.href}
              href={item.href}
              onClick={onNavigate}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors',
                active
                  ? 'bg-sidebar-primary text-sidebar-primary-foreground'
                  : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground',
              )}
            >
              <Icon className="size-5 shrink-0" />
              {item.label}
            </Link>
          )
        })}
      </nav>

      <div className="border-t border-sidebar-border p-3">
        <div className="mb-2 flex items-center gap-3 rounded-lg px-3 py-2">
          <div className="flex size-9 items-center justify-center rounded-full bg-accent text-accent-foreground text-sm font-bold">
            {user?.username?.[0]?.toUpperCase() ?? '؟'}
          </div>
          <div className="min-w-0 leading-tight">
            <p className="truncate text-sm font-medium text-sidebar-foreground">
              {user?.username || 'کاربر'}
            </p>
            <p className="text-xs text-muted-foreground">
              {user?.role ? ROLES[user.role] : '—'}
            </p>
          </div>
        </div>
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 text-muted-foreground hover:text-destructive"
          onClick={logout}
        >
          <LogOut className="size-5" />
          خروج از حساب
        </Button>
      </div>
    </div>
  )
}
