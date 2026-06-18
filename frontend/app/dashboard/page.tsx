'use client'

import Link from 'next/link'
import {
  Package,
  ClipboardList,
  Handshake,
  Users,
  ArrowLeft,
} from 'lucide-react'
import { useAuth } from '@/lib/auth'
import { useApi } from '@/lib/use-api'
import { PageHeader } from '@/components/page-header'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  ORDER_STATUSES,
  ORDER_TYPES,
  ROLES,
  STATUS_VARIANTS,
} from '@/lib/constants'
import type { Order, Product, Partner, User } from '@/lib/types'

export default function DashboardPage() {
  const { user, can } = useAuth()

  const products = useApi<{ products: Product[] }>(
    can('products.read') ? '/products/' : null,
  )
  const orders = useApi<{ orders: Order[] }>(
    can('orders.read') ? '/orders/' : null,
  )
  const partners = useApi<{ partners: Partner[] }>(
    can('partners.read') ? '/partners/' : null,
  )
  const users = useApi<{ users: User[] }>(
    can('users.read') ? '/users/' : null,
  )

  const orderList = orders.data?.orders ?? []
  const pendingCount = orderList.filter((o) => o.status === 'Pending').length

  const stats = [
    {
      show: can('orders.read'),
      label: 'کل سفارش‌ها',
      value: orderList.length,
      loading: orders.loading,
      icon: ClipboardList,
      href: '/dashboard/orders',
      hint: `${pendingCount} در انتظار تأیید`,
    },
    {
      show: can('products.read'),
      label: 'محصولات',
      value: products.data?.products?.length ?? 0,
      loading: products.loading,
      icon: Package,
      href: '/dashboard/products',
    },
    {
      show: can('partners.read'),
      label: 'شرکای تجاری',
      value: partners.data?.partners?.length ?? 0,
      loading: partners.loading,
      icon: Handshake,
      href: '/dashboard/partners',
    },
    {
      show: can('users.read'),
      label: 'کاربران',
      value: users.data?.users?.length ?? 0,
      loading: users.loading,
      icon: Users,
      href: '/dashboard/users',
    },
  ].filter((s) => s.show)

  const recentOrders = orderList.slice(0, 6)

  return (
    <div>
      <PageHeader
        title={`خوش آمدید، ${user?.username ?? ''}`}
        description={`نقش شما: ${user?.role ? ROLES[user.role] : '—'}`}
      />

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {stats.map((stat) => {
          const Icon = stat.icon
          return (
            <Link key={stat.label} href={stat.href}>
              <Card className="transition-shadow hover:shadow-md">
                <CardContent className="flex items-center justify-between gap-4">
                  <div>
                    <p className="text-sm text-muted-foreground">{stat.label}</p>
                    {stat.loading ? (
                      <Skeleton className="mt-2 h-8 w-16" />
                    ) : (
                      <p className="mt-1 text-3xl font-bold text-foreground">
                        {stat.value.toLocaleString('fa-IR')}
                      </p>
                    )}
                    {stat.hint && (
                      <p className="mt-1 text-xs text-muted-foreground">
                        {stat.hint}
                      </p>
                    )}
                  </div>
                  <div className="flex size-12 items-center justify-center rounded-xl bg-accent text-accent-foreground">
                    <Icon className="size-6" />
                  </div>
                </CardContent>
              </Card>
            </Link>
          )
        })}
      </div>

      {can('orders.read') && (
        <Card className="mt-6">
          <CardContent>
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-bold text-foreground">آخرین سفارش‌ها</h2>
              <Link
                href="/dashboard/orders"
                className="flex items-center gap-1 text-sm font-medium text-primary hover:underline"
              >
                مشاهده همه
                <ArrowLeft className="size-4" />
              </Link>
            </div>

            {orders.loading ? (
              <div className="flex flex-col gap-3">
                {Array.from({ length: 4 }).map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            ) : recentOrders.length === 0 ? (
              <p className="py-8 text-center text-sm text-muted-foreground">
                هنوز سفارشی ثبت نشده است.
              </p>
            ) : (
              <ul className="flex flex-col divide-y divide-border">
                {recentOrders.map((order) => (
                  <li
                    key={order.id}
                    className="flex flex-wrap items-center justify-between gap-2 py-3"
                  >
                    <div className="flex items-center gap-3">
                      <span className="font-mono text-sm text-muted-foreground" dir="ltr">
                        {order.order_number}
                      </span>
                      <Badge variant="outline">
                        {ORDER_TYPES[order.order_type] ?? order.order_type}
                      </Badge>
                    </div>
                    <div className="flex items-center gap-3">
                      <span className="text-sm text-foreground">
                        {order.business_partner_name}
                      </span>
                      <Badge variant={STATUS_VARIANTS[order.status] ?? 'secondary'}>
                        {ORDER_STATUSES[order.status] ?? order.status}
                      </Badge>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  )
}
