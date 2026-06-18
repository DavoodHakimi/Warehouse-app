'use client'

import { useState } from 'react'
import { Plus, Pencil, Trash2 } from 'lucide-react'
import { toast } from 'sonner'
import { PageHeader } from '@/components/page-header'
import { DataError, EmptyState, TableSkeleton } from '@/components/data-states'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { PartnerFormDialog } from './partner-form-dialog'
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
import { useApi } from '@/lib/use-api'
import { apiFetch, ApiError } from '@/lib/api'
import { useAuth } from '@/lib/auth'
import { PARTNER_TYPE_LABELS } from '@/lib/constants'
import type { Partner } from '@/lib/types'

export default function PartnersPage() {
  const { can } = useAuth()
  const { data, error, loading, refetch } =
    useApi<{ partners: Partner[] }>('/partners/')
  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Partner | null>(null)
  const [deleting, setDeleting] = useState<Partner | null>(null)

  const partners = data?.partners ?? []
  const canCreate = can('partners.create')
  const canUpdate = can('partners.update')
  const canDelete = can('partners.delete')

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(partner: Partner) {
    setEditing(partner)
    setFormOpen(true)
  }

  async function handleDelete() {
    if (!deleting) return
    try {
      await apiFetch(`/partners/${deleting.id}`, { method: 'DELETE' })
      toast.success('شریک تجاری حذف شد.')
      setDeleting(null)
      refetch()
    } catch (err) {
      toast.error(err instanceof ApiError ? err.message : 'خطا در حذف.')
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <PageHeader
        title="شرکای تجاری"
        description="مدیریت تأمین‌کنندگان و خریداران."
        action={
          canCreate ? (
            <Button onClick={openCreate}>
              <Plus className="size-4" />
              افزودن شریک
            </Button>
          ) : null
        }
      />

      {loading ? (
        <TableSkeleton columns={5} />
      ) : error ? (
        <DataError message={error} onRetry={refetch} />
      ) : partners.length === 0 ? (
        <EmptyState
          title="شریکی یافت نشد"
          description="هنوز شریک تجاری ثبت نشده است."
        />
      ) : (
        <Card className="overflow-hidden p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="text-right">نام شرکت</TableHead>
                <TableHead className="text-right">نوع</TableHead>
                <TableHead className="text-right">شماره تماس</TableHead>
                <TableHead className="text-right">ایمیل</TableHead>
                <TableHead className="text-right">شخص رابط</TableHead>
                {(canUpdate || canDelete) && (
                  <TableHead className="text-left">عملیات</TableHead>
                )}
              </TableRow>
            </TableHeader>
            <TableBody>
              {partners.map((partner) => (
                <TableRow key={partner.id}>
                  <TableCell className="font-medium">{partner.name}</TableCell>
                  <TableCell>
                    <Badge variant="secondary">
                      {PARTNER_TYPE_LABELS[partner.business_partner_type] ??
                        partner.business_partner_type}
                    </Badge>
                  </TableCell>
                  <TableCell dir="ltr" className="text-right">
                    {partner.phone_number}
                  </TableCell>
                  <TableCell dir="ltr" className="text-right">
                    {partner.email}
                  </TableCell>
                  <TableCell>
                    <div className="flex flex-col">
                      <span>{partner.contact_name}</span>
                      <span
                        dir="ltr"
                        className="text-right text-xs text-muted-foreground"
                      >
                        {partner.contact_phone_number}
                      </span>
                    </div>
                  </TableCell>
                  {(canUpdate || canDelete) && (
                    <TableCell className="text-left">
                      <div className="flex justify-start gap-1">
                        {canUpdate && (
                          <Button
                            variant="ghost"
                            size="icon"
                            onClick={() => openEdit(partner)}
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
                            onClick={() => setDeleting(partner)}
                            aria-label="حذف"
                          >
                            <Trash2 className="size-4" />
                          </Button>
                        )}
                      </div>
                    </TableCell>
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Card>
      )}

      <PartnerFormDialog
        open={formOpen}
        onOpenChange={setFormOpen}
        partner={editing}
        onSaved={refetch}
      />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(o) => !o && setDeleting(null)}
        title="حذف شریک تجاری"
        description={`آیا از حذف «${deleting?.name}» مطمئن هستید؟`}
        confirmText="حذف"
        destructive
        onConfirm={handleDelete}
      />
    </div>
  )
}
