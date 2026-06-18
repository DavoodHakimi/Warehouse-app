'use client'

import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { Loader2, Menu, X } from 'lucide-react'
import { useAuth } from '@/lib/auth'
import { DashboardSidebar } from '@/components/dashboard-sidebar'
import { Button } from '@/components/ui/button'

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const { user, loading } = useAuth()
  const router = useRouter()
  const [mobileOpen, setMobileOpen] = useState(false)

  useEffect(() => {
    if (!loading && !user) {
      router.replace('/login')
    }
  }, [user, loading, router])

  if (loading || !user) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="size-8 animate-spin text-primary" />
        <span className="sr-only">در حال بارگذاری</span>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen bg-secondary">
      {/* نوار کناری دسکتاپ */}
      <aside className="fixed inset-y-0 right-0 z-30 hidden w-64 border-l border-border lg:block">
        <DashboardSidebar />
      </aside>

      {/* کشوی موبایل */}
      {mobileOpen && (
        <div className="fixed inset-0 z-40 lg:hidden">
          <div
            className="absolute inset-0 bg-foreground/40"
            onClick={() => setMobileOpen(false)}
            aria-hidden="true"
          />
          <aside className="absolute inset-y-0 right-0 w-64 border-l border-border shadow-xl">
            <Button
              variant="ghost"
              size="icon"
              className="absolute left-2 top-3 z-10"
              onClick={() => setMobileOpen(false)}
              aria-label="بستن منو"
            >
              <X className="size-5" />
            </Button>
            <DashboardSidebar onNavigate={() => setMobileOpen(false)} />
          </aside>
        </div>
      )}

      {/* محتوای اصلی */}
      <div className="flex min-w-0 flex-1 flex-col lg:mr-64">
        <header className="sticky top-0 z-20 flex items-center gap-3 border-b border-border bg-background/80 px-4 py-3 backdrop-blur lg:hidden">
          <Button
            variant="outline"
            size="icon"
            onClick={() => setMobileOpen(true)}
            aria-label="باز کردن منو"
          >
            <Menu className="size-5" />
          </Button>
          <span className="font-bold">سامانه مدیریت انبار</span>
        </header>
        <main className="flex-1 p-4 sm:p-6 lg:p-8">{children}</main>
      </div>
    </div>
  )
}
