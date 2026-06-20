'use client'

import { useState } from 'react'
import { Plus, Trash2, ChevronDown, ArrowRightLeft, Pencil, } from 'lucide-react'
import { toast } from 'sonner'
import { PageHeader } from '@/components/page-header'
import { DataError, EmptyState, TableSkeleton } from '@/components/data-states'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { OrderFormDialog } from './order-form-dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Card } from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useApi } from '@/lib/use-api'
import { apiFetch, ApiError } from '@/lib/api'
import { useAuth } from '@/lib/auth'
import {
  ORDER_STATUSES,
  ORDER_TYPES,
  STATUS_VARIANTS,
} from '@/lib/constants'
import { getTransitions } from '@/lib/order-transitions'
import type { Order } from '@/lib/types'

export default function OrdersPage() {
  const { can } = useAuth()
  const { data, error, loading, refetch } =
    useApi<{ orders: Order[] }>('/orders/')
  const [formOpen, setFormOpen] = useState(false)
  const [formEdit, setFormEdit] = useState<Order | null>(null)
  const [deleting, setDeleting] = useState<Order | null>(null)
  const [pendingAction, setPendingAction] = useState<string | null>(null)

  const orders = data?.orders ?? []
  const canCreate = can('orders.create')
  const canUpdate = can('orders.update')
  const canDelete = can('orders.delete')

  async function runTransition(order: Order, action: string, label: string) {
    setPendingAction(`${order.id}-${action}`)
    try {
      await apiFetch(`/orders/${order.id}/${action}`, { method: 'POST' })
      toast.success(`عملیات «${label}» با موفقیت انجام شد.`)
      refetch()
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در تغییر وضعیت.')
    } finally {
      setPendingAction(null)
    }
  }

  async function handleDelete() {
    if (!deleting) return
    try {
      console.log(deleting.id)
      await apiFetch(`/orders/${deleting.id}`, { method: 'DELETE' })
      toast.success('سفارش حذف شد.')
      setDeleting(null)
      refetch()
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در حذف سفارش.')
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="سفارش‌ها"
        description="مدیریت سفارش‌های خرید و فروش و گردش وضعیت آن‌ها."
        action={
          canCreate ? (
            <Button onClick={() => { setFormOpen(true); setFormEdit(null) }}>
              <Plus className="size-4" />
              ثبت سفارش
            </Button>
          ) : null
        }
      />

      {loading ? (
        <TableSkeleton columns={5} />
      ) : error ? (
        <DataError message={error} onRetry={refetch} />
      ) : orders.length === 0 ? (
        <EmptyState
          title="سفارشی یافت نشد"
          description="هنوز سفارشی ثبت نشده است."
        />
      ) : (
        <Card className="overflow-hidden p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="text-right">شماره سفارش</TableHead>
                <TableHead className="text-right">نوع</TableHead>
                <TableHead className="text-right">شریک تجاری</TableHead>
                <TableHead className="text-right">ارز</TableHead>
                <TableHead className="text-right">وضعیت</TableHead>
                <TableHead className="text-right">عملیات</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {orders.map((order) => {
                const transitions = getTransitions(
                  order.order_type,
                  order.status,
                ).filter((t) => can(t.permission))
                return (
                  <TableRow key={order.id}>
                    <TableCell
                      dir="ltr"
                      className="text-right font-mono text-sm"
                    >
                      {order.order_number}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant={
                          order.order_type === 'sale' ? 'default' : 'outline'
                        }
                      >
                        {ORDER_TYPES[order.order_type] ?? order.order_type}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-medium">
                      {order.business_partner_name}
                    </TableCell>
                    <TableCell>{order.currency}</TableCell>
                    <TableCell>
                      <Badge variant={STATUS_VARIANTS[order.status] ?? 'secondary'}>
                        {ORDER_STATUSES[order.status] ?? order.status}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-left">
                      <div className="flex items-center justify-start gap-1">
                        {transitions.length > 0 && (
                          <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                              <Button variant="outline" size="sm">
                                <ArrowRightLeft className="size-4" />
                                تغییر وضعیت
                                <ChevronDown className="size-4" />
                              </Button>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent align="start">
                              {transitions.map((t) => (
                                <DropdownMenuItem
                                  key={t.action}
                                  disabled={
                                    pendingAction === `${order.id}-${t.action}`
                                  }
                                  onClick={() =>
                                    runTransition(order, t.action, t.label)
                                  }
                                  className={
                                    t.action === 'cancel'
                                      ? 'text-destructive focus:text-destructive'
                                      : ''
                                  }
                                >
                                  {t.label}
                                </DropdownMenuItem>
                              ))}
                            </DropdownMenuContent>
                          </DropdownMenu>
                        )}
                        {canUpdate && (
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => setFormEdit(order)}
                            aria-label="ویرایش"
                          >
                            <Pencil className="size-4" />
                          </Button>
                        )}
                        {canDelete && (
                          <Button
                            variant="ghost"
                            size="icon"
                            className="text-destructive hover:text-destructive"
                            onClick={() => setDeleting(order)}
                            aria-label="حذف"
                          >
                            <Trash2 className="size-4" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  </TableRow>
                )
              })}
            </TableBody>
          </Table>
        </Card>
      )}

      <OrderFormDialog
        open={formOpen || Boolean(formEdit)}
        onOpenChange={(o) => {
          if (!o) {
            setFormOpen(false)
            setFormEdit(null)
          }
        }}
        order={formEdit}
        onSaved={refetch}
      />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(o) => !o && setDeleting(null)}
        title="حذف سفارش"
        description={`آیا از حذف سفارش «${deleting?.order_number}» مطمئن هستید؟`}
        confirmText="حذف"
        destructive
        onConfirm={handleDelete}
      />
    </div>
  )
}
